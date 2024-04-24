-- migrate:up
--
-- Name: VerificationPin; Type: TABLE; Schema: betonz; Owner: -
--
CREATE TABLE "VerificationPin" (
	"pin" text NOT NULL,
	"userId" UUID NOT NULL,
	"createdAt" timestamp(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"updatedAt" timestamp(3) WITH TIME ZONE NOT NULL
);

--
-- Name: VerificationPin VerificationPin_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "VerificationPin"
ADD CONSTRAINT "VerificationPin_pkey" PRIMARY KEY ("pin");

--
-- Name: VerificationPin_userId_key; Type: INDEX; Schema: betonz; Owner: -
--
CREATE UNIQUE INDEX "VerificationPin_userId_key" ON "VerificationPin" USING btree ("userId");

--
-- Name: VerificationPin VerificationPin_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "VerificationPin"
ADD CONSTRAINT "VerificationPin_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" (id) ON UPDATE CASCADE ON DELETE CASCADE;

-- migrate:down
--
-- Name: VerificationPin VerificationPin_userId_fkey; Type: FK CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "VerificationPin" DROP CONSTRAINT "VerificationPin_userId_fkey";

--
-- Name: VerificationPin VerificationPin_userId_key; Type: INDEX; Schema: betonz; Owner: -
--
DROP INDEX IF EXISTS "VerificationPin_userId_key";

--
-- Name: VerificationPin VerificationPin_pkey; Type: CONSTRAINT; Schema: betonz; Owner: -
--
ALTER TABLE ONLY "VerificationPin" DROP CONSTRAINT "VerificationPin_pkey";

--
-- Name: VerificationPin; Type: TABLE; Schema: betonz; Owner: -
--
DROP TABLE IF EXISTS "VerificationPin";
