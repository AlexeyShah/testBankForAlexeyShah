package storage

import "bankService/internal/models"

type Migrator struct {
	tx         *PostgreConnector
	isRollback *bool
}

func NewMigrator() *Migrator {
	conn, err := NewPostgreConnection()
	if err != nil {
		panic(err)
	}

	return &Migrator{
		tx: conn,
	}
}

func (m *Migrator) Magrate() error {
	defer m.close()
	var err error
	defer m.setRollBack(&err)

	err = m.tx.GetTx().AutoMigrate(&models.Wallet{})
	if err != nil {
		return err
	}

	err = m.tx.GetTx().Exec(`
	CREATE OR REPLACE PROCEDURE ballanseUpdate(walletId text, val bigint, op text)
	LANGUAGE plpgsql
	AS $$
	DECLARE current_balance bigint;
	BEGIN
		BEGIN
			IF EXISTS (SELECT 1 FROM wallets WHERE id = walletId) THEN
				SELECT ballance INTO current_balance FROM wallets WHERE id = walletId;

				CASE op
					WHEN 'Deposit' THEN
						UPDATE wallets
						SET ballance = ballance + val
						WHERE id = walletId;
					WHEN 'Withdraw' THEN
						IF current_balance < val THEN
							RAISE EXCEPTION 'Insufficient funds for withdrawal';
						END IF;
						UPDATE wallets
						SET ballance = ballance - val
						WHERE id = walletId;
					ELSE
						RAISE EXCEPTION 'Unknown operation: %', op;
				END CASE;
			ELSE
				CASE op
					WHEN 'Deposit' THEN
						INSERT INTO wallets (id, ballance)
						VALUES (walletId, val);
					WHEN 'Withdraw' THEN
						RAISE EXCEPTION 'Cannot withdraw from a non-existent wallet';
					ELSE
						RAISE EXCEPTION 'Unknown operation: %', op;
				END CASE;
			END IF;
		END;
	END;
	$$;`).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Migrator) close() {
	s.tx.Close(*s.isRollback)
}

func (s *Migrator) setRollBack(err *error) {
	boolVal := false
	if err != nil && *err != nil && len((*err).Error()) > 0 {
		boolVal = true
		s.isRollback = &boolVal
	} else {
		s.isRollback = &boolVal
	}
}
