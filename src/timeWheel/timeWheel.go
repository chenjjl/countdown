package timeWheel

import (
	"countdown/src/event"
	"countdown/src/storage"
	"io/ioutil"
	"time"
)

type TimeWheel struct {
	lilHandTimeWheel *littleHandTimeWheel
	bigHandTimeWheel *bigHandTimeWheel
}

func NewTimeWheel(lhTick time.Duration, lhWheelSize uint64, bhTick time.Duration, bhWheelSize uint64) *TimeWheel {
	lilHandTimeWheel := NewLittleHandTimeWheel(lhTick, lhWheelSize)
	return &TimeWheel{
		lilHandTimeWheel: lilHandTimeWheel,
		bigHandTimeWheel: NewBigHandTimeWheel(bhTick, bhWheelSize, lilHandTimeWheel),
	}
}

func (t *TimeWheel) Start() {
	go t.bigHandTimeWheel.Start()
	go t.lilHandTimeWheel.Start()
	time.Sleep(5 * time.Second) // wait to start
	t.startUp()
}

// startUp load event from file
func (t *TimeWheel) startUp() {
	files, _ := ioutil.ReadDir(storage.DirName)
	if len(files) == 0 {
		return
	}
	log.Infof("files %+v in the dir %s", files, storage.DirName)
	// reload form little hand time wheel file
	needReload, err := storage.ReloadLhEvents()
	if err != nil {
		log.Errorf("failed to reload events")
		log.Error(err)
		return
	}
	for _, e := range needReload {
		err = t.Add(e)
		log.Infof("reload event %+v to little hand time wheel", e)
		if err != nil {
			log.Error(err)
		}
	}
	log.Infof("reload little hand time wheel successful, num of reload events is %d", len(needReload))
	go storage.RemoveAllLhEventFiles()
	// reload from big hand time wheel file
	go func() {
		needReload, err = storage.ReloadBhEvents()
		if err != nil {
			return
		}
		for _, e := range needReload {
			err = t.Add(e)
			log.Infof("reload event %+v to big hand time wheel", e)
			if err != nil {
				log.Error(err)
			}
		}
		log.Infof("reload big hand time wheel successful, num of reload events is %d", len(needReload))
		storage.RemoveBhEventFile()
	}()
}

func (t *TimeWheel) Add(event *event.Event) error {
	if event.Expiration/uint64(time.Second.Milliseconds()) >= t.bigHandTimeWheel.tick {
		err := t.bigHandTimeWheel.Add(event)
		if err != nil {
			log.Error("add event to big hand time wheel fail")
			return err
		}
	} else {
		err := t.lilHandTimeWheel.Add(event)
		if err != nil {
			log.Error("add event to little hand time wheel fail")
			return err
		}
	}
	return nil
}
