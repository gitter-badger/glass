package glass

type Instance interface {
    ProcessSimpleFrame(p *SimpleFrame)
    // Test Frame
    ProcessTestFrame(p *TestFrame)
}
