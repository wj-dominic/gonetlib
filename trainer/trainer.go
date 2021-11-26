package trainer

type Trainer struct {
	id	 uint8

}

func NewTrainer(id uint8) *Trainer{
	return &Trainer{
		id : id,
	}
}