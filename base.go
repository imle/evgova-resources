package resources

import (
	"sync"

	"github.com/imle/resourcefork"
)

//noinspection NonAsciiCharacters
type ResourceLibrary struct {
	Colr *Colr

	Booms map[BoomID]*Boom
	Chars map[CharID]*Char
	Cicns map[CicnID]*Cicn
	Crons map[CronID]*Cron
	Descs map[DescID]*Desc
	Dudes map[DudeID]*Dude
	Flets map[FletID]*Flet
	Govts map[GovtID]*Govt
	Intfs map[IntfID]*Intf
	Junks map[JunkID]*Junk
	Misns map[MisnID]*Misn
	Nebus map[NebuID]*Nebu
	Oopss map[OopsID]*Oops
	Outfs map[OutfID]*Outf
	Perss map[PersID]*Pers
	Picts map[PictID]*Pict
	Ranks map[RankID]*Rank
	RleDs map[RleDID]*RleD
	Roids map[RoidID]*Roid
	Shans map[ShanID]*Shan
	Ships map[ShipID]*Ship
	Snds  map[SndID]*Snd
	Spins map[SpinID]*Spin
	Spobs map[SpobID]*Spob
	StrAs map[StrAID]*StrA
	Systs map[SystID]*Syst
	Weaps map[WeapID]*Weap
}

func NewResourceLibraryFromResourceFork(rf *resourcefork.ResourceFork) *ResourceLibrary {
	rl := &ResourceLibrary{
		Booms: map[BoomID]*Boom{},
		Chars: map[CharID]*Char{},
		Cicns: map[CicnID]*Cicn{},
		Crons: map[CronID]*Cron{},
		Descs: map[DescID]*Desc{},
		Dudes: map[DudeID]*Dude{},
		Flets: map[FletID]*Flet{},
		Govts: map[GovtID]*Govt{},
		Intfs: map[IntfID]*Intf{},
		Junks: map[JunkID]*Junk{},
		Misns: map[MisnID]*Misn{},
		Nebus: map[NebuID]*Nebu{},
		Oopss: map[OopsID]*Oops{},
		Outfs: map[OutfID]*Outf{},
		Perss: map[PersID]*Pers{},
		Picts: map[PictID]*Pict{},
		Ranks: map[RankID]*Rank{},
		RleDs: map[RleDID]*RleD{},
		Roids: map[RoidID]*Roid{},
		Shans: map[ShanID]*Shan{},
		Ships: map[ShipID]*Ship{},
		Snds:  map[SndID]*Snd{},
		Spins: map[SpinID]*Spin{},
		Spobs: map[SpobID]*Spob{},
		StrAs: map[StrAID]*StrA{},
		Systs: map[SystID]*Syst{},
		Weaps: map[WeapID]*Weap{},
	}

	rl.Colr = ColrFromResource(rf.Resources["cölr"][128])

	for id := range rf.Resources["bööm"] {
		rl.Booms[BoomID(id)] = BoomFromResource(rf.Resources["bööm"][id])
	}
	for id := range rf.Resources["chär"] {
		rl.Chars[CharID(id)] = CharFromResource(rf.Resources["chär"][id])
	}
	for id := range rf.Resources["cicn"] {
		rl.Cicns[CicnID(id)] = CicnFromResource(rf.Resources["cicn"][id])
	}
	for id := range rf.Resources["crön"] {
		rl.Crons[CronID(id)] = CronFromResource(rf.Resources["crön"][id])
	}
	for id := range rf.Resources["dësc"] {
		rl.Descs[DescID(id)] = DescFromResource(rf.Resources["dësc"][id])
	}
	for id := range rf.Resources["düde"] {
		rl.Dudes[DudeID(id)] = DudeFromResource(rf.Resources["düde"][id])
	}
	for id := range rf.Resources["flët"] {
		rl.Flets[FletID(id)] = FletFromResource(rf.Resources["flët"][id])
	}
	for id := range rf.Resources["gövt"] {
		rl.Govts[GovtID(id)] = GovtFromResource(rf.Resources["gövt"][id])
	}
	for id := range rf.Resources["ïntf"] {
		rl.Intfs[IntfID(id)] = IntfFromResource(rf.Resources["ïntf"][id])
	}
	for id := range rf.Resources["jünk"] {
		rl.Junks[JunkID(id)] = JunkFromResource(rf.Resources["jünk"][id])
	}
	for id := range rf.Resources["mïsn"] {
		rl.Misns[MisnID(id)] = MisnFromResource(rf.Resources["mïsn"][id])
	}
	for id := range rf.Resources["nëbu"] {
		rl.Nebus[NebuID(id)] = NebuFromResource(rf.Resources["nëbu"][id])
	}
	for id := range rf.Resources["öops"] {
		rl.Oopss[OopsID(id)] = OopsFromResource(rf.Resources["öops"][id])
	}
	for id := range rf.Resources["oütf"] {
		rl.Outfs[OutfID(id)] = OutfFromResource(rf.Resources["oütf"][id])
	}
	for id := range rf.Resources["përs"] {
		rl.Perss[PersID(id)] = PersFromResource(rf.Resources["përs"][id])
	}
	lock := sync.Mutex{}
	//finished := make(chan bool, 10)
	wg := sync.WaitGroup{}
	for id := range rf.Resources["PICT"] {
		wg.Add(1)
		go func(id uint16) {
			temp := PictFromResource(rf.Resources["PICT"][id])

			lock.Lock()
			defer lock.Unlock()
			rl.Picts[PictID(id)] = temp
			wg.Add(-1)
		}(id)
	}
	wg.Wait()
	for id := range rf.Resources["ränk"] {
		rl.Ranks[RankID(id)] = RankFromResource(rf.Resources["ränk"][id])
	}
	lock = sync.Mutex{}
	//finished = make(chan bool, 10)
	wg = sync.WaitGroup{}
	for id := range rf.Resources["rlëD"] {
		wg.Add(1)
		go func(id uint16) {
			temp := RleDFromResource(rf.Resources["rlëD"][id])

			lock.Lock()
			defer lock.Unlock()
			rl.RleDs[RleDID(id)] = temp
			wg.Add(-1)
		}(id)
	}
	wg.Wait()
	for id := range rf.Resources["röid"] {
		rl.Roids[RoidID(id)] = RoidFromResource(rf.Resources["röid"][id])
	}
	for id := range rf.Resources["shän"] {
		rl.Shans[ShanID(id)] = ShanFromResource(rf.Resources["shän"][id])
	}
	for id := range rf.Resources["shïp"] {
		rl.Ships[ShipID(id)] = ShipFromResource(rf.Resources["shïp"][id])
	}
	for id := range rf.Resources["snd"] {
		rl.Snds[SndID(id)] = SndFromResource(rf.Resources["snd"][id])
	}
	for id := range rf.Resources["spïn"] {
		rl.Spins[SpinID(id)] = SpinFromResource(rf.Resources["spïn"][id])
	}
	for id := range rf.Resources["spöb"] {
		rl.Spobs[SpobID(id)] = SpobFromResource(rf.Resources["spöb"][id])
	}
	for id := range rf.Resources["STR#"] {
		rl.StrAs[StrAID(id)] = StrAFromResource(rf.Resources["STR#"][id])
	}
	for id := range rf.Resources["sÿst"] {
		rl.Systs[SystID(id)] = SystFromResource(rf.Resources["sÿst"][id])
	}
	for id := range rf.Resources["wëap"] {
		rl.Weaps[WeapID(id)] = WeapFromResource(rf.Resources["wëap"][id])
	}

	return rl
}
