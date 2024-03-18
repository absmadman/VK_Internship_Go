package server

import "VK_Internship_Go/db"

func AssignmentUsersForUpdateCache(dst *db.User, src db.User) {
	if src.Balance >= 0 {
		dst.Balance = src.Balance
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
}

func (rt *Rout) UpdateUCachesByName(u db.User, name string) {
	uc, ok := rt.uCacheName.Get(name)
	if ok {
		AssignmentUsersForUpdateCache(uc, u)
	}
	rt.uCacheName.Remove(name)
	rt.uCacheName.Add(uc.Name, uc)
	rt.uCacheId.Remove(u.Id)
	rt.uCacheId.Add(u.Id, uc)
}

func (rt *Rout) UpdateUCachesById(u db.User, id int) {
	uc, ok := rt.uCacheId.Get(id)
	if ok {
		AssignmentUsersForUpdateCache(uc, u)
	}
	rt.uCacheId.Remove(id)
	rt.uCacheId.Add(id, uc)
	rt.uCacheName.Remove(uc.Name)
	rt.uCacheName.Add(uc.Name, uc)
}

func AssignmentQuestsForUpdateCache(dst *db.Quest, src db.Quest) {
	if src.Cost >= 0 {
		dst.Cost = src.Cost
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
}

func (rt *Rout) UpdateQCachesByName(q db.Quest, name string) {
	qc, ok := rt.qCacheName.Get(name)
	if ok {
		AssignmentQuestsForUpdateCache(qc, q)
	}
	rt.qCacheName.Remove(name)
	rt.qCacheName.Add(name, qc)
	rt.qCacheId.Remove(q.Id)
	rt.qCacheId.Add(q.Id, qc)
}

func (rt *Rout) UpdateQCachesById(q db.Quest, id int) {
	qc, ok := rt.qCacheId.Get(id)
	if ok {
		AssignmentQuestsForUpdateCache(qc, q)
	}
	rt.qCacheId.Remove(id)
	rt.qCacheId.Add(id, qc)
	rt.qCacheName.Remove(qc.Name)
	rt.qCacheName.Add(qc.Name, qc)
}

func (rt *Rout) RemoveFromUCacheById(id int) {
	uc, ok := rt.uCacheId.Get(id)
	if ok {
		rt.uCacheName.Remove(uc.Name)
		rt.uCacheId.Remove(id)
	}
}

func (rt *Rout) RemoveFromUCacheByName(name string) {
	uc, ok := rt.uCacheName.Get(name)
	if ok {
		rt.uCacheId.Remove(uc.Id)
		rt.uCacheName.Remove(name)
	}
}
