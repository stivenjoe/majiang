package calFanShu

import (
	"majiang-new/card/cardType"
	"majiang-new/common"
)

/*说明：
地胡
*/

const (
	_HUPAI85_ID     = 85
	_HUPAI85_NAME   = "地胡"
	_HUPAI85_FANSHU = 16
	_HUPAI85_KIND   = cardType.HUMETHOD_NORMAL
)

var _HUPAI85_CHECKID_ = []int{} //

func init() {
	fanShuMgr.registerHander(&huPai85{
		id:             _HUPAI85_ID,
		name:           _HUPAI85_NAME,
		fanShu:         _HUPAI85_FANSHU,
		setChcFanShuID: _HUPAI85_CHECKID_,
		huKind:         _HUPAI85_KIND,
	})
}

type huPai85 struct {
	id             int
	name           string
	fanShu         int
	setChcFanShuID []int
	huKind         int
}

func (h *huPai85) GetID() int {
	return h.id
}

func (h *huPai85) Name() string {
	return h.name
}

func (h *huPai85) GetFanShu() int {
	return h.fanShu
}

func (h *huPai85) Satisfy(method *cardType.HuMethod, satisfyedID []int, slBanID []int) (bool, int, []int, []int) {
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

func (h *huPai85) CheckSatisfySelf(method *cardType.HuMethod) bool {
	if method.IsDiHu() {
		return true
	}
	return false
}
