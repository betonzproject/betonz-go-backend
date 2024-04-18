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
		'jun',
		'jun@gmail.com',
		-- jun98765!@#
		'$argon2id$v=19$m=65536,t=3,p=4$fxbpRs9Nn3ndKVH0APXN2Q$E7ggj/mmcnNKI21RK4tG+XWtLYhoYKIQP2egG8xIQZo',
		'demo001',
		'SUPERADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'ee4154df-fa69-45c1-a2a1-ee898bbc4ed9',
		'jinn',
		'jinn@gmail.com',
		-- jinn98765!@$#
		'$argon2id$v=19$m=65536,t=3,p=4$lDNe5VxYEtH1XoFrVTbb8Q$8r3GZ19c+lL3768Cm9M/olgEziPs+uhn8PsmGsNV0zY',
		'demo003',
		'SUPERADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'a2bbb4e7-d7d0-4bfa-923e-5d24224a6f68',
		'meep',
		'squishygasp@gmail.com',
		-- meep98765!@#
		'$argon2id$v=19$m=65536,t=3,p=4$8vEufpgTFL/nIzOd02f5JQ$KEZQSALFYYV2nlAht1NONza7ghgY9DGVYuI0Jp6N414',
		'demo004',
		'SUPERADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'ca631c36-f648-11ee-946c-c7a7f43224db',
		'jinja',
		'jinja@example.com',
		-- jinja98765!@$#
		'$argon2id$v=19$m=65536,t=3,p=4$XBINpGaFEx2T/Tljt3yOPQ$CGnx51ufksnexX2ZkXrwRF85DdnI7f7F7XtM+gVH8lI',
		'demo002',
		'ADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'ef2f83ce-f648-11ee-9dab-3337fcea48ac',
		'ken',
		'ken@example.com',
		-- ken98765!@#
		'$argon2id$v=19$m=65536,t=3,p=4$/0RTWEd0DbpH9feZCXWhUg$pSwI9pmQmAci6qrMN1qTWRGQ74xS74lmAkY4OgUYbtA',
		'demo005',
		'ADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'f99ce41e-f648-11ee-aeca-f757b876c0d8',
		'tom',
		'tom@example.com',
		-- tom98765!@#
		'$argon2id$v=19$m=65536,t=3,p=4$5lO9LdZGkBuj+Y2MI7Z4Pg$Z1Pn1XX+xDV6C/rxiFQzFnfTUdt08PQZUMAUuU2blPY',
		'demo006',
		'ADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'71d4b290-f649-11ee-babb-334a99754b1a',
		'konn',
		'konn@example.com',
		-- konn98765!@$#
		'$argon2id$v=19$m=65536,t=3,p=4$5Rj4OaBzgNwuLWrtZzIfLQ$l02tO7n0OrjLkwOJscCBYCsle/3CkCHBpRQZRz3ZSLs',
		'demo007',
		'ADMIN',
		TRUE,
		now(),
		now()
	),
	(
		'a4fd722e-f649-11ee-b14b-a711ab936b59',
		'jerry',
		'jerry@example.com',
		-- jerry98765!@#
		'$argon2id$v=19$m=65536,t=3,p=4$LvmC8aJ3N+hyCCB7KEI2Ag$n35M/CSZDbojwlOQ6pKXndFVEueTnzwpycgvz0tILl4',
		'demo008',
		'ADMIN',
		TRUE,
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

-- -- Inserting Banks into Bank table
-- INSERT INTO
-- 	"Bank" (id, "userId", name, "accountName", "accountNumber", "createdAt", "updatedAt")
-- VALUES
-- 	(gen_random_uuid (), '413a9f49-a79e-4fc6-b495-2e8c49f469b8', 'AYA', 'Alice', '12345', now(), now()),
-- 	(gen_random_uuid (), 'ee4154df-fa69-45c1-a2a1-ee898bbc4ed9', 'CB', 'Bob', '12346', now(), now()),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'WAVE_PAY',
-- 		'Meep',
-- 		'027680235893218',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'KBZ',
-- 		'Meep',
-- 		'027680235893218',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'AYA',
-- 		'Meep',
-- 		'027680235893218',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'YOMA',
-- 		'Meep',
-- 		'027680235893218',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'OK_DOLLAR',
-- 		'Meep',
-- 		'0988762738',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'KBZPAY',
-- 		'Meep',
-- 		'0988762738',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		gen_random_uuid (),
-- 		'd930b3c5-f837-472c-80b0-78cde6b8220c',
-- 		'KBZPAY',
-- 		'eve',
-- 		'0998712728',
-- 		now(),
-- 		now()
-- 	);
-- -- Inserting TransactionRequests into TransactionRequests table 
-- INSERT INTO
-- 	"TransactionRequest" (
-- 		id,
-- 		"userId",
-- 		"bankName",
-- 		"bankAccountName",
-- 		"bankAccountNumber",
-- 		amount,
-- 		type,
-- 		status,
-- 		"receiptPath",
-- 		"modifiedById",
-- 		"createdAt",
-- 		"updatedAt"
-- 	)
-- VALUES
-- 	(
-- 		1,
-- 		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
-- 		'KBZ',
-- 		'Door Man',
-- 		'889281',
-- 		1000000,
-- 		'WITHDRAW',
-- 		'APPROVED',
-- 		'',
-- 		'4c2a3c9c-5fe0-4581-bc0d-08db1b02c4f9',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		2,
-- 		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
-- 		'KBZ',
-- 		'Door Man',
-- 		'889281',
-- 		1000000,
-- 		'WITHDRAW',
-- 		'PENDING',
-- 		'',
-- 		NULL,
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		3,
-- 		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
-- 		'KBZ',
-- 		'Door Man',
-- 		'889281',
-- 		60000,
-- 		'DEPOSIT',
-- 		'DECLINED',
-- 		'',
-- 		'4c2a3c9c-5fe0-4581-bc0d-08db1b02c4f9',
-- 		now(),
-- 		now()
-- 	),
-- 	(
-- 		4,
-- 		'eae3f8e6-742a-44e2-bf94-c07624e175ed',
-- 		'KBZ',
-- 		'Door Man',
-- 		'889281',
-- 		1000000,
-- 		'DEPOSIT',
-- 		'PENDING',
-- 		'',
-- 		NULL,
-- 		now(),
-- 		now()
-- 	);
