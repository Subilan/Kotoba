const { exec } = require("child_process");

const fs = require("fs");

const kotoba = fs.readFileSync("kotoba.html").toString();

fs.writeFileSync("kotoba.dist.ts", fs.readFileSync("kotoba.ts").toString().replace("[[KOTOBA_TEMPLATE_HTML]]", kotoba));

exec("tsc");

fs.rmSync("kotoba.dist.ts")