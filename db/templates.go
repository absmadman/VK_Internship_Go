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
	selectQNameByName    = `SELECT name FROM quests WHERE name = $1`
	insertQ              = `INSERT INTO quests (name, cost) VALUES ($1, $2)`
	selectQIdByName      = `SELECT id FROM quests WHERE name = $1`
	selectQIdById        = `SELECT id FROM quests WHERE id = $1`
	selectUNameByName    = `SELECT name FROM users WHERE name = $1`
	insertU              = `INSERT INTO users (name, balance) VALUES ($1, $2)`
	selectUIdByName      = `SELECT id FROM users WHERE name = $1`
	selectUIdById        = `SELECT id FROM users WHERE id = $1`
	selectUById          = `SELECT id, name, balance FROM users WHERE id = $1`
	selectUByName        = `SELECT id, name, balance FROM users WHERE name = $1`
	selectQById          = `SELECT id, name, cost FROM quests WHERE id = $1`
	selectQByName        = `SELECT id, name, cost FROM quests WHERE name = $1`
	selectE              = `SELECT user_id, quest_id FROM user_quest WHERE user_id = $1 AND quest_id = $2`
	insertE              = `INSERT INTO user_quest (user_id, quest_id) VALUES ($1, $2)`
	updateBalance        = `UPDATE users SET balance = balance + (SELECT cost FROM quests WHERE id = $1) WHERE id = $2`
	selectBalanceById    = `SELECT balance FROM users WHERE id = $1`
	selectEQIdByUId      = `SELECT quest_id FROM user_quest WHERE user_id = $1`
	deleteQById          = `DELETE FROM quests WHERE id = $1`
	deleteQByName        = `DELETE FROM quests WHERE name = $1`
	deleteUById          = `DELETE FROM users WHERE id = $1`
	deleteUByName        = `DELETE FROM users WHERE name = $1`
)
