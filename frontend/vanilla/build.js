const { exec } = require("child_process");

const fs = require("fs");

const templates_kotoba = fs.readFileSync("src/templates/kotoba.html").toString();
const templates_comment = fs.readFileSync("src/templates/comment.html").toString();

const main_ts = fs.readFileSync("src/main.ts").toString();

fs.writeFileSync("src/dist.ts", main_ts.replace("[[templates.kotoba]]", templates_kotoba).replace("[[templates.comment]]", templates_comment));

exec("npx webpack");