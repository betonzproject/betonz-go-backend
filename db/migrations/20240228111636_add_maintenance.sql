-- migrate:up
ALTER TYPE "EventType"
ADD VALUE 'MAINTENANCE_ADD';

ALTER TYPE "EventType"
ADD VALUE 'MAINTENANCE_UPDATE';

ALTER TYPE "EventType"
ADD VALUE 'MAINTENANCE_DELETE';

--
-- Name: Maintenance; Type: TABLE; Schema: betonz; Owner: -
--
CREATE TABLE "Maintenance" (
	id integer NOT NULL,
	"productCode" integer NOT NULL,
	"maintenancePeriod" tstzrange NOT NULL,
	"gmtOffsetSecs" integer NOT NULL,
	"createdAt" timestamp(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamp(3) WITH TIME ZONE NOT NULL
);

--
-- Name: Maintenance_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--
CREATE SEQUENCE "Maintenance_id_seq" AS integer START
WITH
	1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

--
-- Name: Maintenance_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--
ALTER SEQUENCE "Maintenance_id_seq" OWNED BY "Maintenance".id;

--
-- Name: Maintenance id; Type: DEFAULT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "Maintenance"
ALTER COLUMN id
SET DEFAULT nextval('"Maintenance_id_seq"'::regclass);

--
-- Name: Maintenance Maintenance_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "Maintenance"
ADD CONSTRAINT "Maintenance_pkey" PRIMARY KEY (id);

-- migrate:down
CREATE TYPE "EventType_new" AS ENUM(
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
	'FLAG',
	'SYSTEM_BANK_ADD',
	'SYSTEM_BANK_UPDATE',
	'SYSTEM_BANK_DELETE'
);

ALTER TABLE "Event"
ALTER COLUMN "type" TYPE "EventType_new" USING ("type"::TEXT::"EventType_new");

ALTER TYPE "EventType"
RENAME TO "EventType_old";

ALTER TYPE "EventType_new"
RENAME TO "EventType";

DROP TYPE "EventType_old";

--
-- Name: Maintenance; Type: TABLE; Schema: betonz; Owner: -
--
-- Drop the "Maintenance" table
DROP TABLE IF EXISTS "Maintenance";

--
-- Name: Maintenance_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--
-- Drop the sequence
DROP SEQUENCE IF EXISTS "Maintenance_id_seq";