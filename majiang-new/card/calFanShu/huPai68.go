package calFanShu

import (
	"majiang-new/card/cardType"
	"majiang-new/common"
)

/*说明：
断幺：和牌中没有一、九及字牌。
*/

const (
	_HUPAI68_ID     = 68
	_HUPAI68_NAME   = "断幺"
	_HUPAI68_FANSHU = 2
	_HUPAI68_KIND   = cardType.HUMETHOD_NORMAL
)

var _HUPAI68_CHECKID_ = []int{} //

func init() {
	fanShuMgr.registerHander(&huPai68{
		id:             _HUPAI68_ID,
		name:           _HUPAI68_NAME,
		fanShu:         _HUPAI68_FANSHU,
		setChcFanShuID: _HUPAI68_CHECKID_,
		huKind:         _HUPAI68_KIND,
	})
}

type huPai68 struct {
	id             int
	name           string
	fanShu         int
	setChcFanShuID []int
	huKind         int
}

func (h *huPai68) GetID() int {
	return h.id
}

func (h *huPai68) Name() string {
	return h.name
}

func (h *huPai68) GetFanShu() int {
	return h.fanShu
}

func (h *huPai68) Satisfy(method *cardType.HuMethod, satisfyedID []int, slBanID []int) (bool, int, []int, []int) {
	if method.GetHuPaiKind() != h.huKind {
		return false, 0, satisfyedID, slBanID
	}

	if common.InIntSlace(satisfyedID, h.GetID()) {
		return false, 0, satisfyedID, slBanID
	}

	//不能计算的直接退出
	if common.InIntSlace(slBanID, h.GetID()) {
		return false, 0, satisfyedID, slBanID
	}

	if !h.CheckSatisfySelf(method) {
		slBanID = append(slBanID, h.GetID())
		return false, 0, satisfyedID, slBanID
	}
	//满足后把自己自己要ban的id加入进去
	for _, id := range h.setChcFanShuID {
		if !common.InIntSlace(slBanID, id) {
			slBanID = append(slBanID, id)
		}
	}

	fanShu := h.GetFanShu()
	satisfyedID = append(satisfyedID, h.GetID())
	//再把其他的所有的id全部遍历，有就加上去
	otherChkHander := fanShuMgr.getHanderExcept(append(satisfyedID, slBanID...))
	for _, hander := range otherChkHander {
		ok, tmpFanShu, tmpSatisfyID, slTmpBanID := hander.Satisfy(method, satisfyedID, slBanID)
		slBanID = slTmpBanID
		if ok {
			fanShu += tmpFanShu
			satisfyedID = tmpSatisfyID
		}
	}

	return true, fanShu, satisfyedID, slBanID
}

func (h *huPai68) CheckSatisfySelf(method *cardType.HuMethod) bool {
	ownerIncCard := method.GetAllInclude()
	slWanZi, slTongZi, slTiaoZi, slWord, _ := ownerIncCard.GetTypeSet()
	if slWord != nil {
		return false
	}
	slChk := [...][]uint8{slWanZi, slTongZi, slTiaoZi}
	for _, slChkCard := range slChk {
		if len(slChkCard) == 0 {
			continue
		}
		for _, card := range slChkCard {
			base := card % 10
			if base == 1 || base == 9 {
				return false
			}
		}
	}
	return true
}
