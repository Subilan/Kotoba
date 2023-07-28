interface KotobaConfig {
    commentTextareaPlaceholder: string;
    commentTextareaMaxlength: number;
    commentTextareaMinlength: number;
    commentSameIPLimit: number;
    commentAllowGuest: boolean;
    commentShowGuestIP: boolean;
    commentShowIPOrigin: boolean;
}

export default {
    // 占位符：评论区输入框内显示的默认文本
	commentTextareaPlaceholder: '在此键入评论...',
    // 最大长度：字符数小于或等于这个数字才能发布；默认值 500；设置为 0 无限制。
	commentTextareaMaxlength: 500,
    // 最小长度：字符数大于或等于这个数字才能发布；默认值 0
	commentTextareaMinlength: 0,
    // 同 IP 限额：同一个 IP 所能发表的评论数量；默认值 1
    commentSameIPLimit: 1,
    // 是否允许不登录评论；默认值 false
    commentAllowGuest: false,
    // 【仅针对游客】是否显示评论的 IP；默认值 false
    // 注：此项启用后将会自动向游客显示 IP 会被展示的提示信息
    commentShowGuestIP: false,
    // 是否显示评论的 IP 属地；默认值 false
    commentShowIPOrigin: false,
} as KotobaConfig;
