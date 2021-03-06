package calFanShu

import (
	"majiang-new/card/cardType"
	"majiang-new/common"
)

/*说明：
全求人：全靠吃牌、碰牌、单钓别人打出的牌和牌。不计单钓。
*/

const (
	_HUPAI52_ID     = 52
	_HUPAI52_NAME   = "全求人"
	_HUPAI52_FANSHU = 6
	_HUPAI52_KIND   = cardType.HUMETHOD_NORMAL
)

var _HUPAI52_CHECKID_ = []int{79} //

func init() {
	fanShuMgr.registerHander(&huPai52{
		id:             _HUPAI52_ID,
		name:           _HUPAI52_NAME,
		fanShu:         _HUPAI52_FANSHU,
		setChcFanShuID: _HUPAI52_CHECKID_,
		huKind:         _HUPAI52_KIND,
	})
}

type huPai52 struct {
	id             int
	name           string
	fanShu         int
	setChcFanShuID []int
	huKind         int
}

func (h *huPai52) GetID() int {
	return h.id
}

func (h *huPai52) Name() string {
	return h.name
}

func (h *huPai52) GetFanShu() int {
	return h.fanShu
}

func (h *huPai52) Satisfy(method *cardType.HuMethod, satisfyedID []int, slBanID []int) (bool, int, []int, []int) {
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

func (h *huPai52) CheckSatisfySelf(method *cardType.HuMethod) bool {
	handCard := method.GetAllCard()
	if handCard.Len() != 2 {
		return false
	}
	if method.GetHiddenGangCard() != nil {
		return false
	}
	if method.IsZiMo() {
		return false
	}
	return true
}
