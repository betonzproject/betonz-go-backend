-- migrate:up
DROP TABLE "Report";
DROP TABLE "SessionToken";
DROP TABLE "Transaction";


-- migrate:down
--
-- Name: Report; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE "Report" (
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

CREATE SEQUENCE "Report_id_seq"
	AS integer
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	NO MAXVALUE
	CACHE 1;


--
-- Name: Report_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE "Report_id_seq" OWNED BY "Report".id;


--
-- Name: Report id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Report" ALTER COLUMN id SET DEFAULT nextval('"Report_id_seq"'::regclass);


--
-- Name: Report Report_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Report"
	ADD CONSTRAINT "Report_pkey" PRIMARY KEY (id);


--
-- Name: SessionToken; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE "SessionToken" (
	"tokenHash" text NOT NULL,
	"userId" uuid NOT NULL,
	"createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamp(3) with time zone NOT NULL,
	"expiresAt" timestamp(3) with time zone
);


--
-- Name: SessionToken SessionToken_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "SessionToken"
	ADD CONSTRAINT "SessionToken_pkey" PRIMARY KEY ("tokenHash");


--
-- Name: SessionToken SessionToken_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "SessionToken"
	ADD CONSTRAINT "SessionToken_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: Transaction; Type: TABLE; Schema: betonz; Owner: -
--

CREATE TABLE "Transaction" (
	"initiatorId" uuid NOT NULL,
	"beneficiaryId" uuid NOT NULL,
	product text NOT NULL,
	"balanceBefore" numeric(32,2) NOT NULL,
	"balanceAfter" numeric(32,2) NOT NULL,
	amount numeric(32,2) NOT NULL,
	type "TransactionType" NOT NULL,
	remarks text,
	"createdAt" timestamp(3) with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamp(3) with time zone NOT NULL,
	id integer NOT NULL,
	"receiptPath" text,
	bonus numeric(32,2) DEFAULT 0 NOT NULL
);


--
-- Name: Transaction_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--

CREATE SEQUENCE "Transaction_id_seq"
	AS integer
	START WITH 1
	INCREMENT BY 1
	NO MINVALUE
	NO MAXVALUE
	CACHE 1;


--
-- Name: Transaction_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--

ALTER SEQUENCE "Transaction_id_seq" OWNED BY "Transaction".id;


--
-- Name: Transaction id; Type: DEFAULT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Transaction" ALTER COLUMN id SET DEFAULT nextval('"Transaction_id_seq"'::regclass);


--
-- Name: Transaction Transaction_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Transaction"
	ADD CONSTRAINT "Transaction_pkey" PRIMARY KEY (id);

--
-- Name: Transaction Transaction_beneficiaryId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Transaction"
	ADD CONSTRAINT "Transaction_beneficiaryId_fkey" FOREIGN KEY ("beneficiaryId") REFERENCES "User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Transaction Transaction_initiatorId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--

ALTER TABLE ONLY "Transaction"
	ADD CONSTRAINT "Transaction_initiatorId_fkey" FOREIGN KEY ("initiatorId") REFERENCES "User"(id) ON UPDATE CASCADE ON DELETE RESTRICT;
