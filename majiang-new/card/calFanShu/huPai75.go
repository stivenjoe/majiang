package calFanShu

import (
	"majiang-new/card/cardType"
	"majiang-new/common"
)

/*说明：
缺一门：和牌中缺少一种花色序数牌。
*/

const (
	_HUPAI75_ID     = 75
	_HUPAI75_NAME   = "缺一门"
	_HUPAI75_FANSHU = 1
	_HUPAI75_KIND   = cardType.HUMETHOD_NORMAL
)

var _HUPAI75_CHECKID_ = []int{} //

func init() {
	fanShuMgr.registerHander(&huPai75{
		id:             _HUPAI75_ID,
		name:           _HUPAI75_NAME,
		fanShu:         _HUPAI75_FANSHU,
		setChcFanShuID: _HUPAI75_CHECKID_,
		huKind:         _HUPAI75_KIND,
	})
}

type huPai75 struct {
	id             int
	name           string
	fanShu         int
	setChcFanShuID []int
	huKind         int
}

func (h *huPai75) GetID() int {
	return h.id
}

func (h *huPai75) Name() string {
	return h.name
}

func (h *huPai75) GetFanShu() int {
	return h.fanShu
}

func (h *huPai75) Satisfy(method *cardType.HuMethod, satisfyedID []int, slBanID []int) (bool, int, []int, []int) {
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

func (h *huPai75) CheckSatisfySelf(method *cardType.HuMethod) bool {
	ownerIncAll := method.GetAllInclude()
	if ownerIncAll.GetColorCnt() != 2 {
		return false
	}
	return true
}