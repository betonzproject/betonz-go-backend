-- Inserting user into User table
INSERT INTO
	"User" (
		id,
		username,
		email,
		"passwordHash",
		"etgUsername",
		ROLE,
		"isEmailVerified",
		"updatedAt",
		"createdAt"
	)
VALUES
	(
		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
		'doorman',
		'doorman@example.com',
		-- doorman1
		'$argon2id$v=19$m=65536,t=3,p=4$DJr2f9gqmO/wZeFkr3W0kg$8CQpwB3VP4AM05hhCy7wG6XpENzYdQZ9lT8XXd8QBGU',
		'demo001',
		'PLAYER',
		TRUE,
		now(),
		now()
	),
	(
		'413a9f49-a79e-4fc6-b495-2e8c49f469b8',
		'alice',
		'alice@example.com',
		-- alice111
		'$argon2id$v=19$m=65536,t=3,p=4$+iekk+0k5m6SXFNNfH8XVA$+vPzZqU5BIKt03Pze0t7qelKZKXQ93jgYdR3fKjYnIY',
		'demo002',
		'PLAYER',
		TRUE,
		now(),
		now()
	),
	(
		'ee4154df-fa69-45c1-a2a1-ee898bbc4ed9',
		'bob',
		'bob@example.com',
		-- bob11111
		'$argon2id$v=19$m=65536,t=3,p=4$zVS7kbFFcq+Z9NZIYH/fdg$5Lpf4LiRL4l06aelQSbpIJgsYxDyhC5tyfh9HsN+4/o',
		'demo003',
		'PLAYER',
		TRUE,
		now(),
		now()
	),
	(
		'a2bbb4e7-d7d0-4bfa-923e-5d24224a6f68',
		'meep',
		'squishygasp@gmail.com',
		-- meep1111
		'$argon2id$v=19$m=65536,t=3,p=4$4nyNXia3Bbj/wzb4vRLTaA$x3GApCnmotJWPbyi5Wa2LCzwF6lxMH12qmhwTSj124o',
		'demo004',
		'SUPERADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'4c2a3c9c-5fe0-4581-bc0d-08db1b02c4f9',
		'eve',
		'eve@example.com',
		'$argon2id$v=19$m=65536,t=3,p=4$GRa4BIkOXTW1LhppLIqhbQ$eMLnJO8ecqPTHP/pupgGNbH6NxYEkAwOpJ86d5pgC/g',
		'demo020',
		'ADMIN',
		FALSE,
		now(),
		now()
	),
	(
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'system',
		'system@example.com',
		'$argon2id$v=19$m=16,t=2,p=1$amFqbnNlbmdhdW5n$k4O4SPYtRH6MJsuCLgBErw',
		'system',
		'SYSTEM',
		FALSE,
		now(),
		now()
	);

-- Inserting Banks into Bank table
INSERT INTO
	"Bank" (id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt")
VALUES
	(gen_random_uuid (), '413a9f49-a79e-4fc6-b495-2e8c49f469b8', 'AYA', 'Alice', '12345', now(), now()),
	(gen_random_uuid (), 'ee4154df-fa69-45c1-a2a1-ee898bbc4ed9', 'CB', 'Bob', '12346', now(), now()),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'WAVE_PAY',
		'Meep',
		'027680235893218',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'KBZ',
		'Meep',
		'027680235893218',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'AYA',
		'Meep',
		'027680235893218',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'YOMA',
		'Meep',
		'027680235893218',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'OK_DOLLAR',
		'Meep',
		'0988762738',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'KBZPAY',
		'Meep',
		'0988762738',
		now(),
		now()
	),
	(
		gen_random_uuid (),
		'd930b3c5-f837-472c-80b0-78cde6b8220c',
		'KBZPAY',
		'eve',
		'0998712728',
		now(),
		now()
	);

-- Inserting TransactionRequests into TransactionRequests table 
INSERT INTO
	"TransactionRequest" (
		id,
		"userId",
		"bankName",
		"bankAccountName",
		"bankAccountNumber",
		amount,
		type,
		status,
		"receiptPath",
		"modifiedById",
		"createdAt",
		"updatedAt"
	)
VALUES
	(
		1,
		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
		'KBZ',
		'Door Man',
		'889281',
		1000000,
		'WITHDRAW',
		'APPROVED',
		'',
		'4c2a3c9c-5fe0-4581-bc0d-08db1b02c4f9',
		now(),
		now()
	),
	(
		2,
		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
		'KBZ',
		'Door Man',
		'889281',
		1000000,
		'WITHDRAW',
		'PENDING',
		'',
		NULL,
		now(),
		now()
	),
	(
		3,
		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
		'KBZ',
		'Door Man',
		'889281',
		60000,
		'DEPOSIT',
		'DECLINED',
		'',
		'4c2a3c9c-5fe0-4581-bc0d-08db1b02c4f9',
		now(),
		now()
	),
	(
		4,
		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
		'KBZ',
		'Door Man',
		'889281',
		1000000,
		'DEPOSIT',
		'PENDING',
		'',
		NULL,
		now(),
		now()
	);
