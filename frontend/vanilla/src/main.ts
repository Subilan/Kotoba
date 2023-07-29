import markdownit from "markdown-it";
import {
  KotobaCommentReaction,
  ReactionResultMap,
  KotobaComment,
  KotobaUserPayload,
} from "./types";
import { c, get, getPayloadFromToken } from "./fn";
import { Local } from "./class/local";

const apiURL = c("apiURL");

async function buildReactions(commentUid: string) {
  const reactions = await get<KotobaCommentReaction[]>(
    `${apiURL}/public/get/comment-reactions?uid=${commentUid}`
  );

  if (Array.isArray(reactions)) {
    const resultMap: ReactionResultMap = {};
    for (let r of reactions) {
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
  } else {
    return `<div class="add-reaction">添加贴纸</div>`;
  }
}

(async () => {
  const token = Local.getToken();
  const md = new markdownit();
  let payload: KotobaUserPayload = {
    username: "Guest",
    avatar: "<default>",
    website: window.location.href,
  };
  let template = `[[templates.kotoba]]`;

  if (!token) console.warn("Token not found.");
  else payload = getPayloadFromToken(token);

  template = template
    .replace("macro:user-avatar", payload.avatar)
    .replace("macro:input-placeholder", c("commentTextareaPlaceholder"))
    .replace("macro:input-maxlength", c("commentTextareaMaxlength"))
    .replace("macro:input-minlength", c("commentTextareaMinlength"));

  let commentsHTML = ``;

  const comments = await get<KotobaComment[]>(`${apiURL}/public/get/comments?limit=5&order=desc`);

  if (Array.isArray(comments)) {
    const commentTempl = `[[templates.comment]]`;
    for (let comment of comments) {
      commentsHTML += commentTempl
        .replace("macro:comment-username", comment.username)
        .replace("macro:comment-user-website", comment.user_website || "#")
        .replace("macro:comment-avatar", comment.user_avatar)
        .replace(
          "macro:comment-date",
          comment.updated_at === comment.created_at
            ? new Date(comment.created_at).toLocaleString()
            : `编辑于 ${new Date(comment.updated_at).toLocaleString()}`
        )
        .replace("macro:comment-content-html", md.render(comment.text))
        .replace("macro:reaction-html", await buildReactions(comment.uid));
    }
  }

  template = template.replace("macro:comments-html", commentsHTML);

  document.getElementById("kotoba").innerText = template;
})();
