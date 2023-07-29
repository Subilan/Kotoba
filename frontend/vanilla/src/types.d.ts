import config from "../config";

interface KotobaAccountBasic {
  avatar: string;
  website: string;
}

interface KotobaUserPayload extends KotobaAccountBasic {
  username: string;
}

interface KotobaComment {
  username: string;
  user_avatar: string;
  user_website: string;
  text: string;
  uid: string;
  created_at: number;
  updated_at: number;
}

interface KotobaCommentReaction {
  comment_id: string;
  created_at: number;
  emoji: string;
  username: string;
}

type ReactionResultMap = {
  [prop: string]: {
    count: number;
    meta: {
      created_at: number;
      username: string;
    }[];
  };
};

type ConfigName = keyof typeof config;

type KotobaResponse<T = any> = {
  code: 200 | 500 | 403 | 400 | 404 | 409;
  data: T;
  message: string;
  ok: boolean;
  timestamp: number;
};
