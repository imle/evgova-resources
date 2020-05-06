package resources

import (
	"encoding/binary"

	"github.com/imle/resourcefork"
)

// Cron resources are used to define time-dependent events that occur in a manner that is invisible to the player
// but can cause interesting things to happen in the universe, via the manipulation of control bits.
// Some notes:
// 1. Setting any of the above date fields to 0 or -1 effectively makes that field a wildcard field, which will match to anything.
// 2. If you want an event with a wide possible date range to be guaranteed to never run more than once, make it set a control bit in its OnEnd script that will prevent it from subsequently being eligible for activation.
// 3. The 'M' and 'N' control bit set string operators should probably not be used in conjunction with cron events, unless you really want to confuse the player by moving him around at seemingly random times.
// 4. Local news always takes precedence over independent news, even if there is no corresponding news string to display (the STR# ID must still be greater than zero to not be ignored). You can use this to make everyone in the universe except a particular government or set of governments report on something, for example.

type CronID IDType

type CronFlags struct {
	ContinuousEntry bool // 0x001 - Continuous, iterative cron entry - keep evaluating the cron's OnStart field until the EnableOn expression is no longer true or the constraints of the Require fields are no longer met. This can create infinite loops, so be careful!
	ContinuousExit  bool // 0x002 - Continuous, iterative cron exit - keep evaluating the cron's OnExit field until the EnableOn expression is no longer true or the constraints of the Require fields are no longer met. This can create infinite loops, so be careful!
}

type Cron struct {
	ID CronID

	FirstDay   int16 // The first day of the month (1-31) on which the cron event can be activated. If you set this to 0 or -1, this field will be ignored and only FirstMonth and FirstYear will be considered.
	FirstMonth int16 // The first month of the year (1-12) on which the cron event can be activated. Set to 0 or -1 for this to be ignored.
	FirstYear  int16 // The first year in which the cron event can be activated. Set to 0 or -1 for this to be ignored.

	LastDay   int16 // The last day of the month (1-31) on which the cron event can be activated. Set to 0 or -1 for this to be ignored.
	LastMonth int16 // The last month of the year (1-12) on which the cron event can be activated. Set to 0 or -1 for this to be ignored.
	LastYear  int16 // The last year in which the cron event can be activated. Set to 0 or -1 for this to be ignored.

	Random int16 // The percent chance that the cron event will be activated during the date range defined above. Set to 100 for the event to be activated as soon as it can be.

	Duration    int16 // The duration during which the event is active, in days. If this is set to zero, the event will start and end on the same day, i.e. its OnStart and OnEnd scripts will be run at the same time.
	PreHoldoff  int16 // The number of days to "hold" the event in a waiting state after it is activated and before it starts. Set this to zero to have the event start immediately when it is activated.
	PostHoldoff int16 // The number of days to hold the event in a waiting state after it ends and before it is deactivated. This is used to keep a repeating event from being activated immediately after it has just happened. Set this to zero to have the event be deactivated immediately after it ends.

	Flags    CronFlags
	EnableOn ControlBitTest // A control bit test string that is used to determine whether the cron event is eligible to be activated or not. Leave this blank if you are creating an event whose activation doesn't depend on the state of any control bits.

	OnStart ControlBitFunction // A control bit set string that is called when the cron event starts, after waiting through the PreHoldoff time, if any.
	OnEnd   ControlBitFunction // A control bit set string that is called when the cron event ends.

	Contribute FlagMask64 // When the cron event is active, these two Contribute fields together form a 64-bit flag that is subsequently combined with the Contribute fields from the player's ship and the other outfit items in the player's possession, to be used with the Require fields in the outf and misn resources.
	Require    FlagMask64 // These two Require fields together form a 64-bit flag that is logically and'ed with the Contribute fields from the player's current ship and outfit items. Unless for each 1 bit in the Require fields there is a matching 1 bit in one or more of the Contribute fields, the cron will not be activated. Leave these set to zero if unused.

	NewsGovt    [4]StrAID // On planets or stations that are allied with the government whose ID is given by one of the NewsGovt fields, a string will be randomly selected from the STR# resource whose ID is given by the corresponding GovtNewsStr field, and will be displayed as news while the cron event is active. This allows you to let up to four different governments (and their allies) have their own "local news" for a given cron event. Set unused NewsGovt and GovtNewsStr fields to -1.
	GovtNewsStr [4]StrAID

	IndNewsStr StrAID // The ID of a STR# resource from which to randomly select a string to be displayed in the news dialog while this cron event is in progress, if it doesn't have any applicable local news. Set to -1 for no independent news.
}

func CronFromResource(resource resourcefork.Resource) *Cron {
	return CronFromBytes(CronID(resource.ID), resource.Data)
}

func CronFromBytes(id CronID, b []byte) *Cron {
	flags := binary.BigEndian.Uint16(b[22:])

	t := &Cron{
		ID:          id,
		FirstDay:    int16(binary.BigEndian.Uint16(b[0:])),
		FirstMonth:  int16(binary.BigEndian.Uint16(b[2:])),
		FirstYear:   int16(binary.BigEndian.Uint16(b[4:])),
		LastDay:     int16(binary.BigEndian.Uint16(b[6:])),
		LastMonth:   int16(binary.BigEndian.Uint16(b[8:])),
		LastYear:    int16(binary.BigEndian.Uint16(b[10:])),
		Random:      int16(binary.BigEndian.Uint16(b[12:])),
		Duration:    int16(binary.BigEndian.Uint16(b[14:])),
		PreHoldoff:  int16(binary.BigEndian.Uint16(b[16:])),
		PostHoldoff: int16(binary.BigEndian.Uint16(b[18:])),
		Flags: CronFlags{
			ContinuousEntry: flags&0x0001 == 0x0001,
			ContinuousExit:  flags&0x0002 == 0x0002,
		},
		EnableOn:   ControlBitTest(byteString(b[24:], 254)),
		OnStart:    ControlBitFunction(byteString(b[279:], 255)),
		OnEnd:      ControlBitFunction(byteString(b[534:], 255)),
		Contribute: FlagMask64(binary.BigEndian.Uint64(b[790:])),
		Require:    FlagMask64(binary.BigEndian.Uint64(b[798:])),
		NewsGovt: [4]StrAID{
			StrAID(binary.BigEndian.Uint16(b[806:])),
			StrAID(binary.BigEndian.Uint16(b[808:])),
			StrAID(binary.BigEndian.Uint16(b[810:])),
			StrAID(binary.BigEndian.Uint16(b[812:])),
		},
		GovtNewsStr: [4]StrAID{
			StrAID(binary.BigEndian.Uint16(b[814:])),
			StrAID(binary.BigEndian.Uint16(b[816:])),
			StrAID(binary.BigEndian.Uint16(b[818:])),
			StrAID(binary.BigEndian.Uint16(b[820:])),
		},
		IndNewsStr: StrAID(binary.BigEndian.Uint16(b[20:])),
	}

	return t
}
