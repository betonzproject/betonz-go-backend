-- migrate:up

--
-- Name: betonz; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA betonz;


--
-- Name: BankName; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."BankName" AS ENUM (
    'AGD',
    'AYA',
    'CB',
    'KBZ',
    'KBZPAY',
    'OK_DOLLAR',
    'WAVE_PAY',
    'YOMA'
);


--
-- Name: EventResult; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."EventResult" AS ENUM (
    'SUCCESS',
    'FAIL'
);


--
-- Name: EventType; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."EventType" AS ENUM (
    'LOGIN',
    'REGISTER',
    'PASSWORD_RESET_REQUEST',
    'PASSWORD_RESET_TOKEN_VERIFICATION',
    'PASSWORD_RESET',
    'AUTHENTICATION',
    'AUTHORIZATION',
    'PROFILE_UPDATE',
    'USERNAME_CHANGE',
    'PASSWORD_CHANGE',
    'BANK_ADD',
    'BANK_UPDATE',
    'BANK_DELETE',
    'CHANGE_USER_STATUS',
    'EMAIL_VERIFICATION',
    'ACTIVE',
    'TRANSFER_WALLET',
    'RESTORE_WALLET',
    'TRANSACTION',
    'FLAG'
);


--
-- Name: FlagStatus; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."FlagStatus" AS ENUM (
    'PENDING',
    'RESOLVED',
    'RESTRICTED'
);


--
-- Name: IdentityVerificationStatus; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."IdentityVerificationStatus" AS ENUM (
    'VERIFIED',
    'REJECTED',
    'PENDING'
);


--
-- Name: NotificationType; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."NotificationType" AS ENUM (
    'TRANSACTION',
    'IDENTITY_VERIFICATION'
);


--
-- Name: PromotionType; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."PromotionType" AS ENUM (
    'INACTIVE_BONUS',
    'FIVE_PERCENT_UNLIMITED_BONUS',
    'TEN_PERCENT_UNLIMITED_BONUS'
);


--
-- Name: Role; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."Role" AS ENUM (
    'PLAYER',
    'ADMIN',
    'SUPERADMIN',
    'SYSTEM'
);


--
-- Name: TransactionStatus; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."TransactionStatus" AS ENUM (
    'PENDING',
    'APPROVED',
    'DECLINED'
);


--
-- Name: TransactionType; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."TransactionType" AS ENUM (
    'DEPOSIT',
    'WITHDRAW'
);


--
-- Name: UserStatus; Type: TYPE; Schema: betonz; Owner: -
--

CREATE TYPE betonz."UserStatus" AS ENUM (
    'NORMAL',
    'RESTRICTED'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Bank; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Bank" (
    id uuid NOT NULL,
    "userId" uuid NOT NULL,
    name betonz."BankName" NOT NULL,
    "accountName" text NOT NULL,
    "accountNumber" text NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    disabled boolean DEFAULT false NOT NULL
);


--
-- Name: Bet; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Bet" (
    id integer NOT NULL,
    "refId" text NOT NULL,
    "etgUsername" text NOT NULL,
    "providerUsername" text NOT NULL,
    "productCode" integer NOT NULL,
    "productType" integer NOT NULL,
    "gameId" text,
    details text NOT NULL,
    turnover numeric(32,2) NOT NULL,
    bet numeric(32,2) NOT NULL,
    payout numeric(32,2) NOT NULL,
    status integer NOT NULL,
    "startTime" timestamp(3) with time zone NOT NULL,
    "matchTime" timestamp(3) with time zone NOT NULL,
    "endTime" timestamp(3) with time zone NOT NULL,
    "settleTime" timestamp(3) with time zone NOT NULL,
    "progShare" numeric(32,2) NOT NULL,
    "progWin" numeric(32,2) NOT NULL,
    commission numeric(32,2) NOT NULL,
    "winLoss" numeric(32,2) NOT NULL
);


--
-- Name: Event; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Event" (
    id integer NOT NULL,
    "sourceIp" text,
    "userId" uuid,
    type betonz."EventType" NOT NULL,
    result betonz."EventResult" NOT NULL,
    reason text,
    data jsonb,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    "httpRequest" jsonb
);


--
-- Name: Event_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."Event_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Event_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."Event_id_seq" OWNED BY betonz."Event".id;


--
-- Name: Flag; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Flag" (
    "userId" uuid NOT NULL,
    "modifiedById" uuid,
    reason text,
    remarks text,
    status betonz."FlagStatus" NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL
);


--
-- Name: IdentityVerificationRequests; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."IdentityVerificationRequests" (
    id integer NOT NULL,
    "userId" uuid NOT NULL,
    "modifiedById" uuid,
    status betonz."IdentityVerificationStatus" NOT NULL,
    remarks text,
    "nricFront" text NOT NULL,
    "nricBack" text NOT NULL,
    "holderFace" text NOT NULL,
    "nricName" text NOT NULL,
    nric text NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    dob date
);


--
-- Name: IdentityVerificationRequests_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."IdentityVerificationRequests_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: IdentityVerificationRequests_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."IdentityVerificationRequests_id_seq" OWNED BY betonz."IdentityVerificationRequests".id;


--
-- Name: Notification; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Notification" (
    id integer NOT NULL,
    "userId" uuid NOT NULL,
    type betonz."NotificationType" NOT NULL,
    message text,
    variables jsonb,
    read boolean DEFAULT false NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL
);


--
-- Name: Notification_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."Notification_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Notification_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."Notification_id_seq" OWNED BY betonz."Notification".id;


--
-- Name: PasswordResetToken; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."PasswordResetToken" (
    "tokenHash" text NOT NULL,
    "userId" uuid NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL
);


--
-- Name: Report; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Report" (
    id integer NOT NULL,
    "depositAmount" numeric(65,30) NOT NULL,
    "withdrawAmount" numeric(65,30) NOT NULL,
    "depositCount" numeric(65,30) NOT NULL,
    "withdrawCount" numeric(65,30) NOT NULL,
    "withdrawBankFees" numeric(65,30) NOT NULL,
    "bonusGiven" numeric(65,30) NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "activePlayerCount" numeric(65,30) NOT NULL,
    "inactivePlayerCount" numeric(65,30) NOT NULL,
    "winLoss" numeric(65,30) NOT NULL
);


--
-- Name: Report_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."Report_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Report_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."Report_id_seq" OWNED BY betonz."Report".id;


--
-- Name: SessionToken; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."SessionToken" (
    "tokenHash" text NOT NULL,
    "userId" uuid NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    "expiresAt" timestamp(3) with time zone
);


--
-- Name: Transaction; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."Transaction" (
    "initiatorId" uuid NOT NULL,
    "beneficiaryId" uuid NOT NULL,
    product text NOT NULL,
    "balanceBefore" numeric(32,2) NOT NULL,
    "balanceAfter" numeric(32,2) NOT NULL,
    amount numeric(32,2) NOT NULL,
    type betonz."TransactionType" NOT NULL,
    remarks text,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    id integer NOT NULL,
    "receiptPath" text,
    bonus numeric(32,2) DEFAULT 0 NOT NULL
);


--
-- Name: TransactionRequest; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."TransactionRequest" (
    id integer NOT NULL,
    "userId" uuid NOT NULL,
    "modifiedById" uuid,
    "bankName" betonz."BankName" NOT NULL,
    "bankAccountName" text NOT NULL,
    "bankAccountNumber" text NOT NULL,
    "beneficiaryBankAccountName" text,
    "beneficiaryBankAccountNumber" text,
    amount numeric(32,2) NOT NULL,
    type betonz."TransactionType" NOT NULL,
    "receiptPath" text,
    status betonz."TransactionStatus" NOT NULL,
    remarks text,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    bonus numeric(32,2) DEFAULT 0 NOT NULL,
    "withdrawBankFees" numeric(32,2) DEFAULT 0 NOT NULL,
    "depositToWallet" integer,
    promotion betonz."PromotionType"
);


--
-- Name: TransactionRequest_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."TransactionRequest_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: TransactionRequest_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."TransactionRequest_id_seq" OWNED BY betonz."TransactionRequest".id;


--
-- Name: Transaction_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."Transaction_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: Transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."Transaction_id_seq" OWNED BY betonz."Transaction".id;


--
-- Name: TurnoverTarget; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."TurnoverTarget" (
    id integer NOT NULL,
    "userId" uuid NOT NULL,
    "productCode" integer NOT NULL,
    target numeric(32,2) NOT NULL,
    "promoCode" betonz."PromotionType",
    "transactionRequestId" integer NOT NULL,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL
);


--
-- Name: TurnoverTarget_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE betonz."TurnoverTarget_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: TurnoverTarget_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE betonz."TurnoverTarget_id_seq" OWNED BY betonz."TurnoverTarget".id;


--
-- Name: User; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."User" (
    id uuid NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    "passwordHash" text NOT NULL,
    "displayName" text,
    "phoneNumber" text,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    "etgUsername" text NOT NULL,
    role betonz."Role" DEFAULT 'PLAYER'::betonz."Role" NOT NULL,
    "mainWallet" numeric(32,2) DEFAULT 0 NOT NULL,
    "lastUsedBankId" uuid,
    "profileImage" text,
    status betonz."UserStatus" DEFAULT 'NORMAL'::betonz."UserStatus" NOT NULL,
    "lastLoginIp" text,
    "isEmailVerified" boolean DEFAULT false NOT NULL,
    dob date,
    "lastLoginAt" timestamp(3) without time zone,
    "pendingEmail" text
);


--
-- Name: VerificationToken; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz."VerificationToken" (
    "tokenHash" text NOT NULL,
    "userId" uuid,
    "createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) with time zone NOT NULL,
    "registerInfo" jsonb
);


--
-- Name: _prisma_migrations; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz._prisma_migrations (
    id character varying(36) NOT NULL,
    checksum character varying(64) NOT NULL,
    finished_at timestamp with time zone,
    migration_name character varying(255) NOT NULL,
    logs text,
    rolled_back_at timestamp with time zone,
    started_at timestamp with time zone DEFAULT now() NOT NULL,
    applied_steps_count integer DEFAULT 0 NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE betonz.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: Event id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Event" ALTER COLUMN id SET DEFAULT nextval('betonz."Event_id_seq"'::regclass);


--
-- Name: IdentityVerificationRequests id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."IdentityVerificationRequests" ALTER COLUMN id SET DEFAULT nextval('betonz."IdentityVerificationRequests_id_seq"'::regclass);


--
-- Name: Notification id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Notification" ALTER COLUMN id SET DEFAULT nextval('betonz."Notification_id_seq"'::regclass);


--
-- Name: Report id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Report" ALTER COLUMN id SET DEFAULT nextval('betonz."Report_id_seq"'::regclass);


--
-- Name: Transaction id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Transaction" ALTER COLUMN id SET DEFAULT nextval('betonz."Transaction_id_seq"'::regclass);


--
-- Name: TransactionRequest id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TransactionRequest" ALTER COLUMN id SET DEFAULT nextval('betonz."TransactionRequest_id_seq"'::regclass);


--
-- Name: TurnoverTarget id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TurnoverTarget" ALTER COLUMN id SET DEFAULT nextval('betonz."TurnoverTarget_id_seq"'::regclass);


--
-- Name: Bank Bank_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Bank"
    ADD CONSTRAINT "Bank_pkey" PRIMARY KEY (id);


--
-- Name: Bet Bet_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Bet"
    ADD CONSTRAINT "Bet_pkey" PRIMARY KEY (id);


--
-- Name: Event Event_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Event"
    ADD CONSTRAINT "Event_pkey" PRIMARY KEY (id);


--
-- Name: Flag Flag_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Flag"
    ADD CONSTRAINT "Flag_pkey" PRIMARY KEY ("userId");


--
-- Name: IdentityVerificationRequests IdentityVerificationRequests_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."IdentityVerificationRequests"
    ADD CONSTRAINT "IdentityVerificationRequests_pkey" PRIMARY KEY (id);


--
-- Name: Notification Notification_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Notification"
    ADD CONSTRAINT "Notification_pkey" PRIMARY KEY (id);


--
-- Name: PasswordResetToken PasswordResetToken_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."PasswordResetToken"
    ADD CONSTRAINT "PasswordResetToken_pkey" PRIMARY KEY ("tokenHash");


--
-- Name: Report Report_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Report"
    ADD CONSTRAINT "Report_pkey" PRIMARY KEY (id);


--
-- Name: SessionToken SessionToken_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."SessionToken"
    ADD CONSTRAINT "SessionToken_pkey" PRIMARY KEY ("tokenHash");


--
-- Name: TransactionRequest TransactionRequest_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TransactionRequest"
    ADD CONSTRAINT "TransactionRequest_pkey" PRIMARY KEY (id);


--
-- Name: Transaction Transaction_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Transaction"
    ADD CONSTRAINT "Transaction_pkey" PRIMARY KEY (id);


--
-- Name: TurnoverTarget TurnoverTarget_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TurnoverTarget"
    ADD CONSTRAINT "TurnoverTarget_pkey" PRIMARY KEY (id);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (id);


--
-- Name: VerificationToken VerificationToken_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."VerificationToken"
    ADD CONSTRAINT "VerificationToken_pkey" PRIMARY KEY ("tokenHash");


--
-- Name: _prisma_migrations _prisma_migrations_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz._prisma_migrations
    ADD CONSTRAINT _prisma_migrations_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: PasswordResetToken_userId_key; Type: INDEX; Schema: betonz; Owner: -
--

CREATE UNIQUE INDEX "PasswordResetToken_userId_key" ON betonz."PasswordResetToken" USING btree ("userId");


--
-- Name: User_etgUsername_key; Type: INDEX; Schema: betonz; Owner: -
--

CREATE UNIQUE INDEX "User_etgUsername_key" ON betonz."User" USING btree ("etgUsername");


--
-- Name: User_lastUsedBankId_key; Type: INDEX; Schema: betonz; Owner: -
--

CREATE UNIQUE INDEX "User_lastUsedBankId_key" ON betonz."User" USING btree ("lastUsedBankId");


--
-- Name: User_username_key; Type: INDEX; Schema: betonz; Owner: -
--

CREATE UNIQUE INDEX "User_username_key" ON betonz."User" USING btree (username);


--
-- Name: VerificationToken_userId_key; Type: INDEX; Schema: betonz; Owner: -
--

CREATE UNIQUE INDEX "VerificationToken_userId_key" ON betonz."VerificationToken" USING btree ("userId");


--
-- Name: Bank Bank_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Bank"
    ADD CONSTRAINT "Bank_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Bet Bet_etgUsername_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Bet"
    ADD CONSTRAINT "Bet_etgUsername_fkey" FOREIGN KEY ("etgUsername") REFERENCES betonz."User"("etgUsername") ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Event Event_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Event"
    ADD CONSTRAINT "Event_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: Flag Flag_modifiedById_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Flag"
    ADD CONSTRAINT "Flag_modifiedById_fkey" FOREIGN KEY ("modifiedById") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: Flag Flag_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Flag"
    ADD CONSTRAINT "Flag_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: IdentityVerificationRequests IdentityVerificationRequests_modifiedById_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."IdentityVerificationRequests"
    ADD CONSTRAINT "IdentityVerificationRequests_modifiedById_fkey" FOREIGN KEY ("modifiedById") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: IdentityVerificationRequests IdentityVerificationRequests_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."IdentityVerificationRequests"
    ADD CONSTRAINT "IdentityVerificationRequests_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Notification Notification_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Notification"
    ADD CONSTRAINT "Notification_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: PasswordResetToken PasswordResetToken_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."PasswordResetToken"
    ADD CONSTRAINT "PasswordResetToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SessionToken SessionToken_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."SessionToken"
    ADD CONSTRAINT "SessionToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: TransactionRequest TransactionRequest_modifiedById_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TransactionRequest"
    ADD CONSTRAINT "TransactionRequest_modifiedById_fkey" FOREIGN KEY ("modifiedById") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: TransactionRequest TransactionRequest_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TransactionRequest"
    ADD CONSTRAINT "TransactionRequest_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Transaction Transaction_beneficiaryId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Transaction"
    ADD CONSTRAINT "Transaction_beneficiaryId_fkey" FOREIGN KEY ("beneficiaryId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Transaction Transaction_initiatorId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."Transaction"
    ADD CONSTRAINT "Transaction_initiatorId_fkey" FOREIGN KEY ("initiatorId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: TurnoverTarget TurnoverTarget_transactionRequestId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TurnoverTarget"
    ADD CONSTRAINT "TurnoverTarget_transactionRequestId_fkey" FOREIGN KEY ("transactionRequestId") REFERENCES betonz."TransactionRequest"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: TurnoverTarget TurnoverTarget_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."TurnoverTarget"
    ADD CONSTRAINT "TurnoverTarget_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: User User_lastUsedBankId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."User"
    ADD CONSTRAINT "User_lastUsedBankId_fkey" FOREIGN KEY ("lastUsedBankId") REFERENCES betonz."Bank"(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: VerificationToken VerificationToken_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY betonz."VerificationToken"
    ADD CONSTRAINT "VerificationToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES betonz."User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--



-- migrate:down
DROP SCHEMA betonz CASCADE;
