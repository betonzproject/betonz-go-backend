-- migrate:up
ALTER TYPE "EventType"
ADD VALUE 'REWARD_CLAIM';

ALTER TABLE "User"
ADD COLUMN "betonPoint" numeric(32, 2) DEFAULT 0 NOT NULL;

--
-- Name: InventoryItemType; Type: TYPE; Schema: betonz; Owner: -
--
CREATE TYPE "InventoryItemType" AS ENUM(
	'BONUS',
	'TOKEN_A',
	'TOKEN_B',
	'BETON_POINT',
	'RED_PACK',
	'ROYAL_RED_PACK',
	'RAFFLE_TICKET'
);

--
-- Name: Inventory; Type: TABLE; Schema: betonz; Owner: -
--
CREATE TABLE "Inventory" (
	id integer NOT NULL,
	"userId" uuid NOT NULL,
	item "InventoryItemType" NOT NULL,
	count numeric(32, 2) DEFAULT 0 NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY ("userId") REFERENCES "User" (id),
	CONSTRAINT unique_user_item UNIQUE ("userId", item)
);

--
-- Name: Inventory_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--
CREATE SEQUENCE "Inventory_id_seq" AS integer START
WITH
	1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

--
-- Name: Inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: betonz; Owner: -
--
ALTER SEQUENCE "Inventory_id_seq" OWNED BY "Inventory".id;

--
-- Name: Inventory id; Type: DEFAULT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "Inventory"
ALTER COLUMN id
SET DEFAULT nextval('"Inventory_id_seq"'::regclass);

--
-- Name: Inventory Inventory_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "Inventory"
ADD CONSTRAINT "Inventory_pkey" PRIMARY KEY (id);

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
	'SYSTEM_BANK_DELETE',
	'MAINTENANCE_ADD',
	'MAINTENANCE_UPDATE',
	'MAINTENANCE_DELETE'
);

ALTER TABLE "Event"
ALTER COLUMN "type" TYPE "EventType_new" USING ("type"::TEXT::"EventType_new");

ALTER TYPE "EventType"
RENAME TO "EventType_old";

ALTER TYPE "EventType_new"
RENAME TO "EventType";

DROP TYPE "EventType_old";

--
-- Name: Inventory; Type: TABLE; Schema: betonz; Owner: -
--
-- Drop the "Inventory" table
DROP TABLE IF EXISTS "Inventory";

--
-- Name: Inventory_id_seq; Type: SEQUENCE; Schema: betonz; Owner: -
--
-- Drop the sequence
DROP SEQUENCE IF EXISTS "Inventory_id_seq";

-- Drop the ENUM Type
DROP TYPE "InventoryItemType";

ALTER TABLE "User"
DROP COLUMN "betonPoint";