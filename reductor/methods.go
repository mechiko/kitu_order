package reductor

// ChanIn exposes a send-only handle to the input stream.
// Do not close this channel; the Reductor manages its lifecycle.
func (rdc *Reductor) ChanIn() chan<- Message {
	return rdc.in
}
