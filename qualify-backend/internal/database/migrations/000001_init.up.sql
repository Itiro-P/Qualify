--
-- PostgreSQL database dump
--

-- Dumped from database version 17.9
-- Dumped by pg_dump version 17.9

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', 'public', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: skill; Type: TABLE; Schema: public; Owner: gouser
--

CREATE TABLE public.skill (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.skill OWNER TO gouser;

--
-- Name: skill_id_seq; Type: SEQUENCE; Schema: public; Owner: gouser
--

CREATE SEQUENCE public.skill_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.skill_id_seq OWNER TO gouser;

--
-- Name: skill_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gouser
--

ALTER SEQUENCE public.skill_id_seq OWNED BY public.skill.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: gouser
--

CREATE TABLE public."user" (
    id integer NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    phone character varying(20),
    CONSTRAINT email_format CHECK ((email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'::text))
);


ALTER TABLE public."user" OWNER TO gouser;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: gouser
--

CREATE SEQUENCE public.user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_id_seq OWNER TO gouser;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gouser
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: skill id; Type: DEFAULT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public.skill ALTER COLUMN id SET DEFAULT nextval('public.skill_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Name: skill skill_name_key; Type: CONSTRAINT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public.skill
    ADD CONSTRAINT skill_name_key UNIQUE (name);


--
-- Name: skill skill_pkey; Type: CONSTRAINT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public.skill
    ADD CONSTRAINT skill_pkey PRIMARY KEY (id);


--
-- Name: user user_email_key; Type: CONSTRAINT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_key UNIQUE (email);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: gouser
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--