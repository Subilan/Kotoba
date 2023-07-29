import config from "../../config";
import { KotobaUserPayload, ConfigName } from "./types";

export function getPayloadFromToken(token: string): KotobaUserPayload {
  return JSON.parse(Buffer.from(token.split(".")[1], "base64").toString());
}

export function c(configName: ConfigName) {
    const lim = new Limit(config[configName]);
    switch (configName) {
      case "commentTextareaPlaceholder":
        return lim.default("在此输入评论").mustType("string").collect();
      case "commentTextareaMaxlength":
        return lim.default(500).mustType("number").collect();
      case "commentTextareaMinlength":
        return lim.default(0).mustType("number").collect();
      case "commentSameIPLimit":
        return lim.default(1).mustType("number").collect();
      case "commentAllowGuest":
        return lim.default(false).mustType("boolean").collect();
      case "commentShowGuestIP":
        return lim.default(false).mustType("boolean").collect();
      case "commentShowIPOrigin":
        return lim.default(false).mustType("boolean").collect();
    }
  }