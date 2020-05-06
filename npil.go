package resources

import (
	"io/ioutil"

	"github.com/imle/resourcefork"
)

type MissionObjective struct {
	// Not sure what this looks like
}
type MissionData struct {
	// Not sure what this looks like
}

type NpiL struct {
	LastStellar      SpobID
	ShipClass        ShipID
	Cargo            [6]int16
	_                int16
	Fuel             int16
	Month            int16
	Day              int16
	Year             int16
	Exploration      [2048]int16 // SpobID -> (0: unexplored, 1: visited, 2: landed)
	ItemCount        [512]int16
	LegalStatus      [2048]int16
	WeaponCount      [256]int16
	Ammo             [256]int16
	Cash             Credits
	MissionObjective [16]MissionObjective
	MissionData      [16]MissionData
	MissionBits      [10000]int8
	StellarDominated [2048]int16 // SpobID
	EscortClass      [64]ShipID  // (-1: No Escort, 0-767: Captured, 1000-1767: Hired)
	FighterClass     [64]ShipID  // (-1: No Fighter, 0-767: Fighter)
	_                [128]int16
	UnknownEscorts   [64]ShipID
	CombatRating     int16

	VersionInfo      int16
	StrictPlayFlag   int16
	Gender           int16
	StellarShipCount [2048]int16
	PersonAlive      [1024]int16
	PersonGrudge     [1024]int16
	_                [64]int16
	StellarAnnoyance [2048]int16
	SeenIntroScreen  int8
	_                int8
	DisasterTime     [256]int16
	DisasterStellar  [256]SpobID
	JunkQuantity     [128]int16
	PriceFluctuation [2][2]int16
	CronDuration     [512]int16
	CronHoldOff      [512]int16
	StellarOwned     [2048]int16
	StellarDestroyed [2048]int16
	_                [4]int16
	NicknameLength   int8
	NickName         [63]byte
	ShipColorRed     uint16
	ShipColorGreen   uint16
	ShipColorBlue    uint16
	RankActive       [128]int16
	DatePrefix       [16]byte
	DateSuffix       [16]byte
	_                [1024]int16
}

func PilotFromResource(resource resourcefork.Resource) *NpiL {
	return PilotFromBytes(resource.Data)
}

const EncryptionKey uint32 = 0xb36a210f

func PilotFromFile(path string) (*NpiL, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return PilotFromBytes(data), nil
}

func PilotFromBytes(b []byte) *NpiL {
	t := &NpiL{
		ShipClass: 128,
	}

	return t
}

// int _DoEncryption(int arg0, int arg1, int arg2) {
//    var_12 = arg2;
//    ebx = arg0;
//    ecx = arg1;
//    if (((ecx != 0x0 ? 0x1 : 0x0) & (ebx != 0x0 ? 0x1 : 0x0)) != 0x0) {
//            esi = 0x0;
//            edi = ecx >> 0x2;
//            var_E = ecx - edi * 0x4;
//            ecx = var_12;
//            while (esi != edi) {
//                    esi = esi + 0x1;
//                    edx = ecx >> 0x18 | ecx >> 0x8 & 0xff00 | (ecx & 0xff00) << 0x8 | ecx << 0x18;
//                    ecx = ecx - 0x21524111 ^ 0xdeadbeef;
//                    *ebx = *ebx ^ edx;
//                    ebx = ebx + 0x4;
//            }
//            eax = ebx;
//            edx = ecx >> 0x18 | ecx >> 0x8 & 0xff00 | (ecx & 0xff00) << 0x8 | ecx << 0x18;
//            ecx = 0x0;
//            edx = edx ^ *ebx;
//            ebx = -(var_E & 0xffff);
//            do {
//                    eax = eax + 0x1;
//                    if (ecx == ebx) {
//                        break;
//                    }
//                    *(int8_t *)(eax - 0x1) = edx;
//                    ecx = ecx - 0x1;
//                    edx = edx >> 0x8;
//            } while (true);
//    }
//    return 0x0;
//}
func encrypt(b []uint32, key uint32) []byte {
	if len(b) == 0 {
		return nil
	}

	var eax uint32 = 0
	var ebx uint32 = 0
	var ecx uint32 = key
	var edx uint32 = 0

	var esi uint32 = 0x0

	// Rounding thing
	var edi uint32 = uint32(len(b)) >> 0x2 // length divided by 4
	var E uint32 = uint32(len(b)) - edi*0x4

	for esi != edi {
		esi++
		edx = ecx>>0x18 | ecx>>0x8&0xff00 | (ecx&0xff00)<<0x8 | ecx<<0x18
		ecx = ecx - 0x21524111 ^ 0xdeadbeef
		b[ebx] = b[ebx] ^ edx
		ebx = ebx + 0x4
	}

	eax = ebx
	edx = ecx>>0x18 | ecx>>0x8&0xff00 | (ecx&0xff00)<<0x8 | ecx<<0x18
	ecx = 0x0
	edx = edx ^ b[ebx]
	ebx = -(E & 0xffff)

	for {
		eax = eax + 0x1
		if ecx == ebx {
			break
		}
		//*(int8_t *)(eax - 0x1) = edx
		ecx = ecx - 0x1
		edx = edx >> 0x8
	}

	return nil
}

func decrypt(b []byte, key uint32) []byte {
	// int _LoadPilotData(int arg0) {
	//    var_4 = arg0;
	//    eax = *dword_deeb8;
	//    ebx = *dword_dee7c;
	//    *eax = 0x0;
	//    *ebx = 0x0;
	//    eax = FSResolveAliasFile(var_4, 0x1, &var_19, &var_1A);
	//    if ((eax != 0x0) || (var_19 != 0x0)) goto loc_7649d;
	//
	//loc_7510c:
	//    eax = _EVOpenResourceFile(var_4, 0x3);
	//    *(int16_t *)*dword_deecc = eax;
	//    eax = eax + 0x1;
	//    if (eax == 0x0) goto loc_7649d;
	//
	//loc_75130:
	//    ebx = *dword_dee7c;
	//    _crashlog("load pilot data 1");
	//    eax = GetResource(0x4e70954c, 0x80);
	//    *ebx = eax;
	//    if (eax == 0x0) goto loc_75bbf;
	//
	//loc_7515a:
	//    ebx = *dword_dee7c;
	//    HNoPurge();
	//    var_484 = GetHandleSize(*ebx);
	//    eax = *ebx;
	//    eax = *eax;
	//    eax = _DoEncryption(eax, var_484, 0xb36a210f);
	//    var_448 = 0x0;
	//    if (eax != 0x0) goto loc_75bc7;
	//
	//loc_75194:
	//    ebx = *dword_dee7c;
	//    eax = *ebx;
	//    GetHandleSize(eax);
	//    _FlipNpil(0x80, **ebx);
	//    _crashlog("load pilot data 2");
	//    _PropagateMissionBitEffects();
	//    _HandleStellarSystemVisibility();
	//    eax = *ebx;
	//    ecx = *0x3b5138;
	//    edx = *(int16_t *)*eax & 0xffff;
	//    var_458 = ecx;
	//    *(int16_t *)(*ecx + 0x92) = edx;
	//    if (edx == 0xffff) goto loc_7521c;
	//
	//loc_751f6:
	//    var_448 = 0x0;
	//    goto loc_7525b;
	//
	//loc_7525b:
	//    ecx = sign_extend_32(edx);
	//    edi = ecx * 0x498;
	//    *(int16_t *)(*var_458 + 0x74) = *(int16_t *)(edi + **0x3b5010 + 0x14) & 0xffff;
	//    ebx = *var_458;
	//    eax = *(int16_t *)(ebx + 0x74) & 0xffff;
	//    if (eax <= 0x7ff) goto loc_752af;
	//
	//loc_75289:
	//    *(int16_t *)(ebx + 0x74) = _FindSystemFromStellar(ecx);
	//    var_458 = *0x3b5138;
	//    goto loc_752ce;
	//
	//loc_752ce:
	//    if (*(int16_t *)(*var_458 + 0x74) <= 0x800) goto loc_7531f;
	//
	//loc_752de:
	//    _crashlog("warning: pilot file system ID out of range");
	//    edx = 0x0;
	//    *(int16_t *)(*var_458 + 0x74) = 0x0;
	//    eax = *0x3b5180;
	//    eax = *eax;
	//    goto loc_75301;
	//
	//loc_75301:
	//    if (*(int8_t *)(eax + 0x1ed) != 0x0) goto loc_752a2;
	//
	//loc_7530a:
	//    edx = edx + 0x1;
	//    eax = eax + 0x1fc;
	//    if (edx != 0x800) goto loc_75301;
	//
	//loc_75318:
	//    var_448 = 0x1;
	//    goto loc_7531f;
	//
	//loc_7531f:
	//    _PropagateMissionBitEffects();
	//    _HandleStellarSystemVisibility();
	//    ecx = *0x3b5010;
	//    esi = *0x3b5138;
	//    eax = *ecx;
	//    edx = *esi;
	//    xmm0 = intrinsic_cvtsi2ss(xmm0, sign_extend_32(*(int16_t *)(edi + eax + 0x4)));
	//    *(edx + 0x18) = intrinsic_movss(*(edx + 0x18), xmm0);
	//    eax = *ecx;
	//    edx = *esi;
	//    xmm0 = intrinsic_cvtsi2ss(xmm0, sign_extend_32(*(int16_t *)(edi + eax + 0x6)));
	//    *(edx + 0x1c) = intrinsic_movss(*(edx + 0x1c), xmm0);
	//    ebx = *esi;
	//    eax = _Rand(0x168);
	//    edx = *0x3b5180;
	//    xmm0 = intrinsic_cvtsi2ss(xmm0, sign_extend_16(eax));
	//    *(ebx + 0x44) = intrinsic_movss(*(ebx + 0x44), xmm0);
	//    ecx = *esi;
	//    edx = *edx;
	//    ebx = *0x3b51b8;
	//    *(int16_t *)ebx = *(int16_t *)(sign_extend_32(*(int16_t *)(ecx + 0x74)) * 0x1fc + edx) & 0xffff;
	//    *(int16_t *)(ebx + 0x2) = *(int16_t *)(sign_extend_32(*(int16_t *)(ecx + 0x74)) * 0x1fc + edx + 0x2) & 0xffff;
	//    eax = **dword_dee7c;
	//    ebx = *0x3b5098;
	//    *(int16_t *)(ecx + 0x76) = *(int16_t *)(*eax + 0x2) & 0xffff;
	//    var_458 = esi;
	//    if (*(int16_t *)(sign_extend_32(*(int16_t *)(*esi + 0x76)) * 0xabc + *ebx + 0xa) != 0xd8f1) goto loc_75436;
	//
	//loc_753da:
	//    esi = *0x3b5138;
	//    ebx = *0x3b5098;
	//    _crashlog("warning: pilot file ship class out of range");
	//    edx = 0x0;
	//    *(int16_t *)(*esi + 0x76) = 0x0;
	//    eax = *ebx;
	//    goto loc_753f2;
	//
	//loc_753f2:
	//    if (*(int16_t *)(eax + 0xa) != 0xd8f1) goto loc_7541d;
	//
	//loc_753fa:
	//    edx = edx + 0x1;
	//    eax = eax + 0xabc;
	//    if (edx != 0x300) goto loc_753f2;
	//
	//loc_75408:
	//    var_448 = 0x1;
	//    var_458 = *0x3b5138;
	//    goto loc_75436;
	//
	//loc_75436:
	//    ecx = 0x0;
	//    do {
	//            *(int16_t *)(*var_458 + ecx * 0x2 + 0x7a) = *(int16_t *)(***dword_dee7c + ecx * 0x2 + 0x4) & 0xffff;
	//            ecx = ecx + 0x1;
	//    } while (ecx != 0x6);
	//    esi = *0x3b51a8;
	//    ebx = 0x0;
	//    edi = 0x0;
	//    do {
	//            eax = **dword_dee7c;
	//            eax = *eax;
	//            eax = *(int16_t *)(eax + ebx * 0x2 + 0x101a) & 0xffff;
	//            *(int16_t *)esi = eax;
	//            if (eax > 0x0) {
	//                    if (*(int16_t *)(edi + **0x3b5194) == 0x7fff) {
	//                            _crashlog_print("warning: pilot file has invalid outfit item ID ");
	//                            _crashlog_print_num(ebx + 0x80);
	//                            _crashlog_print("\n");
	//                            *(int16_t *)esi = 0x0;
	//                            var_448 = 0x1;
	//                    }
	//            }
	//            ebx = ebx + 0x1;
	//            esi = esi + 0x2;
	//            edi = edi + 0x37c;
	//    } while (ebx != 0x200);
	//    ebx = *0x3b5138;
	//    esi = *ebx;
	//    _ShipShieldCapacity(esi);
	//    asm { fstp       dword [esi+0x54] };
	//    esi = *ebx;
	//    _ShipArmorCapacity(esi);
	//    ecx = *dword_dee7c;
	//    asm { fstp       dword [esi+0x58] };
	//    eax = *ecx;
	//    edx = *ebx;
	//    ebx = *dword_dee7c;
	//    ecx = *0x3b5154;
	//    *(edx + 0x38) = intrinsic_movss(*(edx + 0x38), intrinsic_cvtsi2ss(xmm0, sign_extend_32(*(int16_t *)(*eax + 0x12))));
	//    edx = *ebx;
	//    ebx = 0x0;
	//    *(int16_t *)(ecx + 0x2) = *(int16_t *)(*edx + 0x14) & 0xffff;
	//    *(int16_t *)(ecx + 0x4) = *(int16_t *)(*edx + 0x16) & 0xffff;
	//    *(int16_t *)ecx = *(int16_t *)(*edx + 0x18) & 0xffff;
	//    ecx = 0x0;
	//    do {
	//            esi = *0x3b5180;
	//            edi = *dword_dee7c;
	//            edx = *esi;
	//            eax = *edi;
	//            eax = *eax;
	//            eax = *(int16_t *)(eax + ecx * 0x2 + 0x1a) & 0xffff;
	//            ecx = ecx + 0x1;
	//            *(int16_t *)(ebx + edx + 0x90) = eax;
	//            ebx = ebx + 0x1fc;
	//    } while (ecx != 0x800);
	//    edi = *dword_dee7c;
	//    ecx = *0x3b5134;
	//    edx = 0x0;
	//    ebx = *edi;
	//    do {
	//            eax = *ebx;
	//            eax = *(int16_t *)(eax + edx * 0x2 + 0x141a) & 0xffff;
	//            edx = edx + 0x1;
	//            *(int16_t *)ecx = eax;
	//            ecx = ecx + 0x2;
	//    } while (edx != 0x800);
	//    edi = 0x80;
	//    var_42C = 0x0;
	//    var_458 = *0x3b5138;
	//    do {
	//            ebx = edi - 0x80;
	//            esi = ebx * 0xc8;
	//            edx = esi + *var_458;
	//            ecx = *dword_dee7c;
	//            *(int16_t *)(edx + 0xc8) = *(int16_t *)(**ecx + ebx * 0x2 + 0x241a) & 0xffff;
	//            *(int16_t *)(0xd0 + esi + *var_458) = *(int16_t *)(**ecx + ebx * 0x2 + 0x261a) & 0xffff;
	//            if ((*(int16_t *)(0xc8 + esi + *var_458) > 0x0) && (*(int16_t *)(var_42C + **0x3b517c + 0x6) == 0xd8f1)) {
	//                    _crashlog_print("warning: pilot file has invalid weapon ID ");
	//                    _crashlog_print_num(edi);
	//                    _crashlog_print("\n");
	//                    edx = *0x3b5138;
	//                    *(int16_t *)(0xc8 + esi + *var_458) = 0x0;
	//                    var_448 = 0x1;
	//                    var_458 = edx;
	//            }
	//            ebx = ebx * 0xc8;
	//            if (*(int16_t *)(0xd0 + ebx + *var_458) > 0x0) {
	//                    if (*(int16_t *)(var_42C + **0x3b517c + 0x6) == 0xd8f1) {
	//                            _crashlog_print("warning: pilot file has invalid weapon ammo ID ");
	//                            _crashlog_print_num(edi);
	//                            _crashlog_print("\n");
	//                            edx = *0x3b5138;
	//                            *(int16_t *)(0xd0 + ebx + *var_458) = 0x0;
	//                            var_448 = 0x1;
	//                            var_458 = edx;
	//                    }
	//            }
	//            edi = edi + 0x1;
	//            var_42C = var_42C + 0xc8;
	//    } while (edi != 0x180);
	//    ebx = *dword_dee7c;
	//    edi = 0x0;
	//    esi = 0x0;
	//    *(*var_458 + 0xa0) = *(**ebx + 0x281a);
	//    **0x3b512c = *(**ebx + 0xe9ae);
	//    _crashlog("load pilot data 3");
	//    var_44C = ebx;
	//    var_430 = 0x0;
	//    var_450 = &var_220;
	//    var_460 = 0x0;
	//    var_464 = 0x0;
	//    do {
	//            edx = *0x3b50f0;
	//            eax = *var_44C;
	//            ecx = *edx;
	//            eax = *eax;
	//            *(edi + ecx) = *(eax + var_464 + 0x281e);
	//            *(edi + ecx + 0x4) = *(eax + var_464 + 0x2822);
	//            *(edi + ecx + 0x8) = *(eax + var_464 + 0x2826);
	//            *(edi + ecx + 0xc) = *(eax + var_464 + 0x282a);
	//            *(edi + ecx + 0x10) = *(eax + var_464 + 0x282e);
	//            memcpy(esi + **0x3b50f4, 0x295e + **var_44C + var_460, 0x8ec);
	//            edx = edi + **0x3b50f0;
	//            if (*(int8_t *)edx != 0x0) {
	//                    _SetMissionDeadline(edx + 0x6, sign_extend_32(*(int16_t *)(esi + **0x3b50f4 + 0x48)));
	//                    ebx = *0x3b50f4;
	//                    *(int8_t *)(esi + *ebx + 0x70) = 0x0;
	//                    *(int8_t *)(esi + *ebx + 0xb0) = 0x0;
	//                    eax = esi + *ebx;
	//                    edx = *(int16_t *)(eax + 0x4a) & 0xffff;
	//                    if (edx > 0x7f) {
	//                            eax = *(int16_t *)(eax + 0x4c) & 0xffff;
	//                            if (eax > 0x0) {
	//                                    ebx = *0x3b50f4;
	//                                    GetIndString();
	//                                    strncpy(0x70 + esi + *ebx, var_450, 0x3f);
	//                            }
	//                    }
	//                    eax = esi + **0x3b50f4;
	//                    edx = *(int16_t *)(eax + 0x52) & 0xffff;
	//                    if (edx > 0x7f) {
	//                            eax = *(int16_t *)(eax + 0x54) & 0xffff;
	//                            if (eax > 0x0) {
	//                                    var_45C = *0x3b50f4;
	//                                    GetIndString();
	//                                    strncpy(0xb0 + esi + *var_45C, var_450, 0x3f);
	//                            }
	//                    }
	//                    ecx = *0x3b50f4;
	//                    eax = *ecx;
	//                    edx = esi + eax;
	//                    if ((*(int8_t *)(edx + 0x58) & 0x10) != 0x0) {
	//                            *(int16_t *)(edx + 0x6e) = *(int16_t *)(edx + 0x64) & 0xffff;
	//                            eax = *ecx;
	//                    }
	//                    var_45C = *0x3b50f4;
	//                    *(int16_t *)(0x6c + esi + eax) = _Rand(0x46) + 0x46;
	//                    *(int16_t *)(esi + *var_45C + 0x6a) = 0x0;
	//            }
	//            var_430 = var_430 + 0x1;
	//            edi = edi + 0x14;
	//            var_464 = var_464 + 0x14;
	//            esi = esi + 0x8ec;
	//            var_460 = var_460 + 0x8ec;
	//    } while (var_430 != 0x10);
	//    ecx = 0x0;
	//    do {
	//            *(int8_t *)(ecx + **0x3b5220) = *(int8_t *)(ecx + ***dword_dee7c + 0xb81e) & 0xff;
	//            ecx = ecx + 0x1;
	//    } while (ecx != 0x2710);
	//    ecx = 0x0;
	//    ebx = 0x0;
	//    do {
	//            edx = ebx + **0x3b5010;
	//            if (*(int8_t *)(edx + 0x45) != 0x0) {
	//                    *(int8_t *)(edx + 0x46) = *(int8_t *)(ecx + ***dword_dee7c + 0xdf2e) & 0xff;
	//            }
	//            else {
	//                    *(int8_t *)(edx + 0x46) = 0x0;
	//            }
	//            ecx = ecx + 0x1;
	//            ebx = ebx + 0x498;
	//    } while (ecx != 0x800);
	//    esi = 0x0;
	//    _DeleteAllComputerShips(0x1);
	//    do {
	//            ecx = *(int16_t *)(***dword_dee7c + esi * 0x2 + 0xe72e) & 0xffff;
	//            if (ecx >= 0x0) {
	//                    ecx = sign_extend_32(ecx - ((SAR(ecx + (sign_extend_32(ecx) * 0xffff8313 >> 0x10), 0x9)) - (SAR(ecx, 0xf))) * 0x3e8);
	//                    if (*(int16_t *)(ecx * 0xabc + **0x3b5098 + 0xa) == 0xd8f1) {
	//                            _crashlog("warning: pilot file has invalid escort ship type");
	//                            var_448 = 0x1;
	//                    }
	//                    else {
	//                            eax = _SpawnEscort(ecx, 0xffffffff);
	//                            edx = eax;
	//                            if (eax != 0xffff) {
	//                                    if (*(int16_t *)(***dword_dee7c + esi * 0x2 + 0xe72e) <= 0x3e7) {
	//                                            ebx = *0x3b5138;
	//                                            ecx = sign_extend_32(edx);
	//                                            edx = ecx * 0xc940;
	//                                            var_458 = ebx;
	//                                            eax = *ebx;
	//                                            ebx = *dword_dee7c;
	//                                            *(int8_t *)(edx + eax + 0xbb) = 0x0;
	//                                            if (*(int16_t *)(**ebx + esi * 0x2 + 0xe8ae) != 0x0) {
	//                                                    *(int8_t *)(edx + *var_458 + 0xbe) = 0x1;
	//                                            }
	//                                    }
	//                                    else {
	//                                            eax = *0x3b5138;
	//                                            ecx = sign_extend_32(edx);
	//                                            var_458 = eax;
	//                                            *(int8_t *)(ecx * 0xc940 + *eax + 0xbb) = 0x1;
	//                                    }
	//                                    if (*(int16_t *)(***dword_dee7c + esi * 0x2 + 0xe82e) != 0x0) {
	//                                            *(int8_t *)(ecx * 0xc940 + *var_458 + 0xbf) = 0x1;
	//                                    }
	//                            }
	//                    }
	//            }
	//            eax = **dword_dee7c;
	//            eax = *eax;
	//            eax = *(int16_t *)(eax + esi * 0x2 + 0xe7ae) & 0xffff;
	//            if (eax <= 0x2ff) {
	//                    ecx = sign_extend_32(eax);
	//                    if (*(int16_t *)(ecx * 0xabc + **0x3b5098 + 0xa) == 0xd8f1) {
	//                            _crashlog("warning: pilot file has invalid fighter ship type");
	//                            var_448 = 0x1;
	//                    }
	//                    else {
	//                            eax = _SpawnEscort(ecx, 0xffffffff);
	//                            if (eax != 0xffff) {
	//                                    ecx = *0x3b5138;
	//                                    ebx = sign_extend_16(eax) * 0xc940;
	//                                    *(int16_t *)(ebx + *ecx + 0x88) = 0x5;
	//                                    *(int8_t *)(ebx + *ecx + 0xbb) = 0x0;
	//                                    _AISetup(ebx + *ecx);
	//                                    edx = *(int16_t *)(***dword_dee7c + esi * 0x2 + 0xe92e) & 0xffff;
	//                                    if (edx != 0xffff) {
	//                                            *(int16_t *)(ebx + **0x3b5138 + 0xc922) = edx;
	//                                    }
	//                            }
	//                    }
	//            }
	//            esi = esi + 0x1;
	//    } while (esi != 0x40);
	//    ebx = *dword_dee7c;
	//    HPurge();
	//    ReleaseResource(*ebx);
	//    _crashlog("load pilot data 4");
	//    goto loc_75bc7;
	//
	//loc_75bc7:
	//    _crashlog("load pilot data 5");
	//    eax = GetResource(0x4e70954c, 0x81);
	//    **dword_deeb8 = eax;
	//    if (eax == 0x0) goto loc_7622d;
	//
	//loc_75bf7:
	//    ebx = *dword_deeb8;
	//    HNoPurge();
	//    if (_DoEncryption(**ebx, GetHandleSize(*ebx), 0xb36a210f) != 0x0) goto loc_76205;
	//
	//loc_75c2a:
	//    ebx = *dword_deeb8;
	//    eax = *ebx;
	//    GetHandleSize(eax);
	//    _FlipNpil(0x81, **ebx);
	//    _crashlog("load pilot data 6");
	//    eax = *ebx;
	//    ecx = *eax;
	//    edx = *(int16_t *)ecx & 0xffff;
	//    if (edx != 0x6b) goto loc_75c99;
	//
	//loc_75c6d:
	//    HUnlock(eax);
	//    HPurge();
	//    CloseResFile(sign_extend_32(*(int16_t *)*dword_deecc));
	//    eax = 0xffffffd3;
	//    return eax;
	//
	//loc_75c99:
	//    if (edx > 0x12b) goto loc_75ccc;
	//
	//loc_75ca0:
	//    HUnlock(eax);
	//    HPurge();
	//    CloseResFile(sign_extend_32(*(int16_t *)*dword_deecc));
	//    eax = 0xffffffd6;
	//    return eax;
	//
	//loc_75ccc:
	//    if (*(int16_t *)(ecx + 0x2) == 0x1) {
	//            *(int8_t *)*0x3b5158 = 0x1;
	//    }
	//    else {
	//            *(int8_t *)*0x3b5158 = 0x0;
	//    }
	//    esi = *dword_deeb8;
	//    var_454 = esi;
	//    if (*(int16_t *)(**esi + 0x4) == 0x1) {
	//            *(int16_t *)*0x3b5224 = 0x1;
	//    }
	//    else {
	//            *(int16_t *)*0x3b5224 = 0x0;
	//    }
	//    ebx = 0x0;
	//    strncpy(*0x3b5290, **var_454 + 0x5d98, 0x40);
	//    eax = *var_454;
	//    ecx = *0x3b5084;
	//    edx = *eax;
	//    *ecx = *(edx + 0x5dd8);
	//    *(int16_t *)(ecx + 0x4) = *(int16_t *)(edx + 0x5ddc) & 0xffff;
	//    ecx = 0x0;
	//    do {
	//            edx = ecx + **0x3b5010;
	//            if ((*(int8_t *)(edx + 0x45) != 0x0) && (*(int8_t *)(edx + 0x46) == 0x0)) {
	//                    var_454 = *dword_deeb8;
	//                    esi = *0x3b5010;
	//                    *(int16_t *)(edx + 0x466) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x6) & 0xffff;
	//                    *(int16_t *)(ecx + *esi + 0x2a) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x2086) & 0xffff;
	//            }
	//            else {
	//                    esi = *0x3b5010;
	//                    *(int16_t *)(edx + 0x466) = 0x0;
	//                    *(int16_t *)(ecx + *esi + 0x2a) = 0x0;
	//            }
	//            ebx = ebx + 0x1;
	//            ecx = ecx + 0x498;
	//    } while (ebx != 0x800);
	//    ebx = 0x0;
	//    ecx = 0x0;
	//    do {
	//            edx = ecx + **0x3b50bc;
	//            if (*(int16_t *)(edx + 0x4) > 0x0) {
	//                    if (*(int8_t *)(edx + 0x623) != 0x0) {
	//                            edi = *dword_deeb8;
	//                            if (*(int16_t *)(**edi + ebx * 0x2 + 0x1006) != 0x0) {
	//                                    *(int8_t *)(edx + 0x620) = 0x1;
	//                                    if (*(int16_t *)(**edi + ebx * 0x2 + 0x1806) != 0x0) {
	//                                            *(int8_t *)(ecx + **0x3b50bc + 0x621) = 0x1;
	//                                    }
	//                                    else {
	//                                            *(int8_t *)(ecx + **0x3b50bc + 0x621) = 0x0;
	//                                    }
	//                            }
	//                            else {
	//                                    *(int8_t *)(edx + 0x620) = 0x0;
	//                                    *(int8_t *)(ecx + **0x3b50bc + 0x621) = 0x0;
	//                            }
	//                    }
	//                    else {
	//                            *(int8_t *)(edx + 0x620) = 0x0;
	//                            *(int8_t *)(ecx + **0x3b50bc + 0x621) = 0x0;
	//                    }
	//            }
	//            else {
	//                    *(int8_t *)(edx + 0x620) = 0x0;
	//            }
	//            ebx = ebx + 0x1;
	//            ecx = ecx + 0x794;
	//    } while (ebx != 0x400);
	//    ebx = 0x0;
	//    esi = 0x0;
	//    do {
	//            if (*(int16_t *)(***dword_deeb8 + ebx * 0x2 + 0x5dde) == 0x0) {
	//                    *(int8_t *)(esi + **0x3b5090) = 0x0;
	//            }
	//            else {
	//                    eax = esi + **0x3b5090;
	//                    if (*(int8_t *)(eax + 0x1) != 0x0) {
	//                            *(int8_t *)eax = 0x1;
	//                    }
	//                    else {
	//                            edi = *0x3b5090;
	//                            _crashlog("warning: pilot file contains invalid rank");
	//                            *(int8_t *)(*edi + esi) = 0x0;
	//                            var_448 = 0x1;
	//                    }
	//            }
	//            ebx = ebx + 0x1;
	//            esi = esi + 0x120;
	//    } while (ebx != 0x80);
	//    ebx = ebx & 0xffffff00;
	//    _crashlog("load pilot data 7");
	//    ecx = *dword_deeb8;
	//    var_454 = ecx;
	//    eax = *ecx;
	//    ecx = 0x0;
	//    *(int8_t *)*0x3b5260 = *(int8_t *)(*eax + 0x3086) & 0xff;
	//    *(int16_t *)*0x3b5430 = 0xffff;
	//    do {
	//            edx = ecx + **0x3b52ec;
	//            if (*(int8_t *)(edx + 0x20e) != 0x0) {
	//                    esi = *0x3b52ec;
	//                    *(int16_t *)(edx + 0xc) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x3088) & 0xffff;
	//                    *(int16_t *)(ecx + *esi + 0x2) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x3288) & 0xffff;
	//            }
	//            else {
	//                    esi = *0x3b52ec;
	//                    *(int16_t *)(edx + 0xc) = 0x0;
	//                    *(int16_t *)(*esi + ecx + 0x2) = 0xffff;
	//            }
	//            ebx = ebx + 0x1;
	//            ecx = ecx + 0x210;
	//    } while (ebx != 0x100);
	//    esi = 0x0;
	//    ebx = 0x0;
	//    do {
	//            edi = *0x3b50cc;
	//            *(int16_t *)(*edi + ebx + 0x22) = *(int16_t *)(**var_454 + esi * 0x2 + 0x3488) & 0xffff;
	//            eax = ebx + *edi;
	//            if ((*(int16_t *)(eax + 0x22) > 0x0) && (*(int16_t *)(eax + 0x20) < 0x0)) {
	//                    edi = *0x3b50cc;
	//                    _crashlog("warning: pilot file contains invalid junk type");
	//                    *(int16_t *)(*edi + ebx + 0x22) = 0x0;
	//                    var_448 = 0x1;
	//                    var_454 = *dword_deeb8;
	//            }
	//            esi = esi + 0x1;
	//            ebx = ebx + 0x526;
	//    } while (esi != 0x80);
	//    ebx = 0x0;
	//    ecx = 0x0;
	//    do {
	//            edx = ecx + **0x3b5024;
	//            if (*(int8_t *)(edx + 0x1) != 0x0) {
	//                    esi = *0x3b5024;
	//                    *(int16_t *)(edx + 0x26) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x3590) & 0xffff;
	//                    *(int16_t *)(*esi + ecx + 0x28) = *(int16_t *)(**var_454 + ebx * 0x2 + 0x3990) & 0xffff;
	//                    eax = ecx + *esi;
	//                    if ((*(int16_t *)(eax + 0x26) < 0x0) && (*(int16_t *)(eax + 0x28) < 0x0)) {
	//                            *(int8_t *)eax = 0x0;
	//                    }
	//                    else {
	//                            *(int8_t *)eax = 0x1;
	//                    }
	//            }
	//            else {
	//                    *(int8_t *)edx = 0x0;
	//            }
	//            ebx = ebx + 0x1;
	//            ecx = ecx + 0x350;
	//    } while (ebx != 0x200);
	//    ecx = 0x0;
	//    ebx = 0x0;
	//    do {
	//            edx = **0x3b5180;
	//            eax = *var_454;
	//            eax = *eax;
	//            eax = *(int16_t *)(eax + ecx * 0x2 + 0x3d90) & 0xffff;
	//            ecx = ecx + 0x1;
	//            *(int16_t *)(edx + ebx + 0xc4) = eax;
	//            ebx = ebx + 0x1fc;
	//    } while (ecx != 0x800);
	//    ebx = 0x0;
	//    ecx = 0x0;
	//    do {
	//            edx = *(int16_t *)(**var_454 + ebx * 0x2 + 0x4d90) & 0xffff;
	//            if (edx > 0x0) {
	//                    esi = *0x3b5010;
	//                    *(int16_t *)(*esi + ecx + 0x47c) = edx;
	//                    *(*esi + ecx + 0x3c) = 0xffffffff;
	//            }
	//            else {
	//                    esi = *0x3b5010;
	//                    *(int16_t *)(*esi + ecx + 0x47c) = 0xffff;
	//                    edx = ecx + *esi;
	//                    *(edx + 0x3c) = *(edx + 0x40);
	//            }
	//            ebx = ebx + 0x1;
	//            ecx = ecx + 0x498;
	//    } while (ebx != 0x800);
	//    ecx = 0x0;
	//    edi = *0x3b5394;
	//    ebx = *var_454;
	//    edx = edi;
	//    do {
	//            *(int16_t *)edx = 0xffff;
	//            eax = *ebx;
	//            eax = *(int16_t *)(eax + ecx * 0x2 + 0x5d90) & 0xffff;
	//            if (eax >= 0x0) {
	//                    *(int16_t *)edx = eax;
	//            }
	//            ecx = ecx + 0x1;
	//            edx = edx + 0x2;
	//    } while (ecx != 0x4);
	//    ebx = 0xc940;
	//    var_458 = *0x3b5138;
	//    do {
	//            ecx = ebx + *var_458;
	//            if ((*(int8_t *)(ecx + 0xb8) != 0x0) && (*(int16_t *)(ecx + 0x9a) == 0x0)) {
	//                    *(int16_t *)(ecx + 0xc90a) = *(int16_t *)(*0x3b5394 + sign_extend_32(*(int16_t *)(sign_extend_32(*(int16_t *)(ecx + 0x76)) * 0xabc + **0x3b5098 + 0xe)) * 0x2) & 0xffff;
	//            }
	//            ebx = ebx + 0xc940;
	//    } while (ebx != 0x325000);
	//    edx = 0x0;
	//    ebx = *0x3b5370;
	//    ecx = *0x3b53ac;
	//    esi = *var_454;
	//    do {
	//            *(int16_t *)ebx = *(int16_t *)(*esi + edx * 0x2 + 0x3588) & 0xffff;
	//            eax = *esi;
	//            ebx = ebx + 0x2;
	//            eax = *(int16_t *)(eax + edx * 0x2 + 0x358c) & 0xffff;
	//            edx = edx + 0x1;
	//            *(int16_t *)ecx = eax;
	//            ecx = ecx + 0x2;
	//    } while (edx != 0x2);
	//    esi = *0x3b5184;
	//    strncpy(esi, **var_454 + 0x5ede, 0xf);
	//    ebx = *0x3b50c0;
	//    strncpy(ebx, **var_454 + 0x5eee, 0xf);
	//    *(int8_t *)(esi + 0xf) = 0x0;
	//    *(int8_t *)(ebx + 0xf) = 0x0;
	//    HPurge();
	//    ReleaseResource(*var_454);
	//    _crashlog("load pilot data 8");
	//    goto loc_76205;
	//
	//loc_76205:
	//    _crashlog("load pilot data 9");
	//    eax = GetResource(0x4e70954c, 0x81);
	//    esi = eax;
	//    if (eax != 0x0) {
	//            HNoPurge();
	//            GetResInfo(esi, &var_1C, &var_20, &var_220);
	//            strncpy(*0x3b5210, &var_220, 0x40);
	//            HPurge();
	//            ReleaseResource(esi);
	//    }
	//    FSGetCatalogInfo(var_4, 0x0, 0x0, &var_420, 0x0, 0x0);
	//    _UnicodeNameGetHFSName(var_420 & 0xffff, &var_41E, 0xffff, 0x0, &var_120);
	//    strncpy(*0x3b525c, &var_120, 0x1f);
	//    CloseResFile(sign_extend_32(*(int16_t *)*dword_deecc));
	//    _crashlog("making alias");
	//    _MakeAlias(var_4);
	//    edi = *0x3b5138;
	//    var_458 = edi;
	//    *(int16_t *)(*edi + 0xc8fa) = 0x0;
	//    *(*edi + 0xc8e4) = 0x0;
	//    ebx = *edi;
	//    *(int16_t *)(ebx + 0xc8fc) = (*(int16_t *)(sign_extend_32(*(int16_t *)(ebx + 0x76)) * 0xabc + **0x3b5098 + 0xa00) & 0xffff) - 0x1;
	//    *(*var_458 + 0xc8f0) = 0x0;
	//    *(*var_458 + 0xc8e8) = 0x0;
	//    *(*var_458 + 0xc910) = 0x0;
	//    ecx = 0xc940;
	//    do {
	//            ebx = *var_458;
	//            edx = ebx + ecx;
	//            if ((*(int8_t *)(edx + 0xb8) != 0x0) && (*(int16_t *)(edx + 0x9a) == 0x0)) {
	//                    edi = *0x3b5098;
	//                    eax = sign_extend_32(*(int16_t *)(edx + 0x76));
	//                    esi = *edi;
	//                    if ((*(int8_t *)(eax * 0xabc + esi + 0xa24) & 0x80) != 0x0) {
	//                            *(int16_t *)(edx + 0xc8fc) = (*(int16_t *)(sign_extend_32(*(int16_t *)(ebx + 0x76)) * 0xabc + esi + 0xa00) & 0xffff) - 0x1;
	//                    }
	//                    else {
	//                            *(int16_t *)(edx + 0xc8fc) = 0x0;
	//                    }
	//            }
	//            ecx = ecx + 0xc940;
	//    } while (ecx != 0x325000);
	//    esi = 0x0;
	//    do {
	//            edi = *0x3b5098;
	//            *(int16_t *)(0xa2a + esi + *edi) = _Rand(0x64) + 0x1;
	//            ebx = esi;
	//            esi = esi + 0xabc;
	//            *(int16_t *)(0xa2c + ebx + *edi) = _Rand(0x64) + 0x1;
	//    } while (esi != 0x203400);
	//    esi = 0x0;
	//    do {
	//            ebx = esi;
	//            esi = esi + 0x37c;
	//            *(int16_t *)(0x1e + ebx + **0x3b5194) = _Rand(0x64) + 0x1;
	//    } while (esi != 0x6f800);
	//    _crashlog("dumping pilot log");
	//    _DumpPilotLog();
	//    eax = SAR((var_448 & 0xff) << 0x1f, 0x1f) & 0xffffffd2;
	//    return eax;
	//
	//loc_7622d:
	//    eax = *dword_deecc;
	//    eax = sign_extend_32(*(int16_t *)eax);
	//    goto loc_76235;
	//
	//loc_76235:
	//    CloseResFile(eax);
	//    goto loc_7649d;
	//
	//loc_7649d:
	//    eax = 0xffffffd5;
	//    return eax;
	//
	//loc_7541d:
	//    edi = *0x3b5138;
	//    var_458 = edi;
	//    *(int16_t *)(*edi + 0x76) = edx;
	//    var_448 = 0x1;
	//    goto loc_75436;
	//
	//loc_752a2:
	//    *(int16_t *)(**0x3b5138 + 0x74) = edx;
	//    goto loc_75318;
	//
	//loc_752af:
	//    if (*(int8_t *)(sign_extend_16(eax) * 0x1fc + **0x3b5180 + 0x1ed) != 0x0) goto loc_7531f;
	//
	//loc_752c8:
	//    *(int16_t *)(ebx + 0x74) = 0xffff;
	//    goto loc_752ce;
	//
	//loc_7521c:
	//    _crashlog("warning: pilot file stellar ID out of range");
	//    eax = 0x0;
	//    edx = **0x3b5010;
	//    goto loc_75232;
	//
	//loc_75232:
	//    if (*(int8_t *)(edx + 0x45) != 0x0) goto loc_75205;
	//
	//loc_75238:
	//    eax = eax + 0x1;
	//    edx = edx + 0x498;
	//    if (eax != 0x800) goto loc_75232;
	//
	//loc_75246:
	//    edx = 0x0;
	//    var_448 = 0x1;
	//    var_458 = *0x3b5138;
	//    goto loc_7525b;
	//
	//loc_75205:
	//    edx = eax;
	//    var_448 = 0x1;
	//    var_458 = *0x3b5138;
	//    goto loc_7525b;
	//
	//loc_75bbf:
	//    eax = sign_extend_32(*(int16_t *)*dword_deecc);
	//    goto loc_76235;
	//}

	return b
}
