import axios from "node_modules/axios/index";
import markdownit from "markdown-it";
import { KotobaCommentReaction, ReactionResultMap, KotobaComment } from "./types";
import { c, getPayloadFromToken } from "./fn";

const apiURL = c("apiURL");

async function buildReactions(commentUid: string) {
  const reactions = await axios.get<KotobaCommentReaction[]>(
    `${apiURL}/public/get/comment-reactions?uid=${commentUid}`
  );

  if (!reactions.data) {
    return `<div class="add-reaction">添加贴纸</div>`;
  } else {
    const resultMap: ReactionResultMap = {};
    for (let r of reactions.data) {
      if (Object.keys(resultMap).includes(r.emoji)) {
        resultMap[r.emoji].count++;
        resultMap[r.emoji].meta.push({
          created_at: r.created_at,
          username: r.username,
        });
      } else {
        resultMap[r.emoji] = {
          count: 1,
          meta: [
            {
              created_at: r.created_at,
              username: r.username,
            },
          ],
        };
      }
    }
    let resultHTML = `<div class="reactions">`;
    Object.keys(resultMap).forEach((emoji) => {
      resultHTML += `<span data-meta="${JSON.stringify(resultMap[emoji].meta)}">${emoji} · ${
        resultMap[emoji].count
      }</span>`;
    });
    resultHTML += "</div>";
    return resultHTML;
  }
}

(async () => {
  const token = Local.getToken();
  const md = new markdownit();
  let template = `[[templates.kotoba]]`;

  if (!token) {
    console.warn("Token not found.");
  }

  const payload = getPayloadFromToken(token);

  template = template
    .replace("macro:user-avatar", payload.avatar)
    .replace("macro:input-placeholder", c("commentTextareaPlaceholder"))
    .replace("macro:input-maxlength", c("commentTextareaMaxlength"))
    .replace("macro:input-minlength", c("commentTextareaMinlength"));

  const initialCommentsReq = await axios.get<KotobaComment[]>(
    `${apiURL}/public/get/comments?limit=5&order=desc`
  );

  if (initialCommentsReq.data) {
    const initialComments = initialCommentsReq.data;
    const commentTempl = `[[templates.comment]]`;
    let comments = ``;

    for (let comment of initialComments) {
      comments += commentTempl
        .replace("macro.comment-avatar", comment.user_avatar)
        .replace(
          "macro.comment-date",
          comment.updated_at === comment.created_at
            ? new Date(comment.created_at).toLocaleString()
            : `编辑于 ${new Date(comment.updated_at).toLocaleString()}`
        )
        .replace("macro.comment-content-html", md.render(comment.text))
        .replace("macro.reaction-html", await buildReactions(comment.uid));
    }
  }
})();