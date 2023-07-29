import config from "../config";
import { Limit } from "./class/limit";
import { KotobaUserPayload, ConfigName, KotobaResponse } from "./types";

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
    default:
      return config[configName];
  }
}

export async function get<T = any>(url: string): Promise<T> {
  const resp = await window.fetch(url, {
    method: "GET",
  });

  if (resp.ok) {
    const json = await resp.json() as KotobaResponse<T>;
    return json.data;
  } else {
    return null;
  }
}
