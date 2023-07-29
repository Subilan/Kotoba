const { exec } = require("child_process");

const fs = require("fs");

const templates_kotoba = fs.readFileSync("src/templates/kotoba.html").toString();
const templates_comment = fs.readFileSync("src/templates/comment.html").toString();

const main_ts = fs.readFileSync("src/main.ts").toString();

fs.writeFileSync("src/main.ts.backup", main_ts);
fs.writeFileSync("src/main.ts", main_ts.replace("[[templates.kotoba]]", templates_kotoba).replace("[[templates.comment]]", templates_comment));

exec("tsc");

fs.rmSync("src/main.ts")
fs.renameSync("src/main.ts.backup", "src/main.ts");