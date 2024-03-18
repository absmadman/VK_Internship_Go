package db

const (
	updateUBalanceByName = `UPDATE users SET balance = $1 WHERE name = $2`
	updateUBalanceById   = `UPDATE users SET balance = $1 WHERE id = $2`
	updateUNameById      = `UPDATE users SET name = $1 WHERE id = $2`
	updateUNameByName    = `UPDATE users SET name = $1 WHERE name = $2`
	updateQCostById      = `UPDATE quests SET cost = $1 WHERE id = $2`
	updateQNameById      = `UPDATE quests SET name = $1 WHERE id = $2`
	updateQCostByName    = `UPDATE quests SET cost = $1 WHERE name = $2`
	updateQNameByName    = `UPDATE quests SET name = $1 WHERE name = $2`
	selectQByName        = `SELECT name FROM quests WHERE name = $1`
	insertQ              = `INSERT INTO quests (name, cost) VALUES ($1, $2)`
	selectQIdByName      = `SELECT id FROM quests WHERE name = $1`
	selectUByName        = `SELECT name FROM users WHERE name = $1`
	insertU              = `INSERT INTO users (name, balance) VALUES ($1, $2)`
	selectUNameByName    = `SELECT name FROM users WHERE name = $1`
	selectUIdByName      = `SELECT id FROM users WHERE name = $1`
)
