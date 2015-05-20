package controls

type SwitchUpdatedEvent struct {
    swtch Switch
}

func CreateSwitchUpdatedEvent(swtch Switch) SwitchUpdatedEvent {
    return SwitchUpdatedEvent{swtch}
}

func (event SwitchUpdatedEvent) Switch() Switch {
    return event.swtch
}
