package channels

type PushMeChannel struct{ *BaseChannel }

func NewPushMeChannel() Channel {
	return &PushMeChannel{NewBaseChannel(ChannelPushMe, []string{FormatTypeText})}
}

func (c *PushMeChannel) Send(config ChannelConfig, msg *Message) (*Result, error) {
	pushKey := config.GetString("push_key")
	if pushKey == "" {
		return SendError("pushme config missing: push_key is required"), nil
	}

	cli := PushMe{
		PushKey: pushKey,
		URL:     config.GetString("url"),
		Date:    config.GetString("date"),
		Type:    config.GetString("type"),
	}

	res, err := cli.Request(msg.Title, msg.Text)
	if err != nil {
		return ErrorResult(res, err), nil
	}
	return SuccessResult(res), nil
}
