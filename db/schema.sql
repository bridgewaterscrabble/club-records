--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Ubuntu 15.2-1.pgdg22.04+1)
-- Dumped by pg_dump version 15.2 (Ubuntu 15.2-1.pgdg22.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bingos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bingos (
    game integer NOT NULL,
    player integer NOT NULL,
    word character varying(16) NOT NULL,
    score smallint,
    valid character(1) DEFAULT ' '::bpchar,
    challenged boolean DEFAULT false,
    no_blank boolean,
    triple_triple boolean DEFAULT false,
    non_bingo boolean DEFAULT false,
    CONSTRAINT bingos_score_check CHECK ((score > 50))
);


--
-- Name: games; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.games (
    id integer NOT NULL,
    session integer NOT NULL,
    round smallint NOT NULL,
    player_1 integer NOT NULL,
    player_2 integer NOT NULL,
    score_1 smallint NOT NULL,
    score_2 smallint NOT NULL,
    blank_1 character(1) NOT NULL,
    blank_2 character(1) NOT NULL,
    blank_1_p integer,
    blank_2_p integer,
    lex_family character varying(16) DEFAULT 'nwl'::character varying,
    CONSTRAINT games_check CHECK ((player_1 <> player_2)),
    CONSTRAINT games_check1 CHECK (((blank_1_p IS NULL) OR (blank_1_p = player_1) OR (blank_1_p = player_2))),
    CONSTRAINT games_check2 CHECK (((blank_2_p IS NULL) OR (blank_2_p = player_1) OR (blank_2_p = player_2)))
);


--
-- Name: games_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.games_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: games_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.games_id_seq OWNED BY public.games.id;


--
-- Name: lexica; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lexica (
    lexicon character varying(16) NOT NULL,
    lex_family character varying(16) DEFAULT 'nwl'::character varying,
    start_date date DEFAULT '1900-01-01'::date,
    end_date date DEFAULT '2999-12-31'::date,
    id integer NOT NULL
);


--
-- Name: lexica_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.lexica_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: lexica_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.lexica_id_seq OWNED BY public.lexica.id;


--
-- Name: lexicon_families; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.lexicon_families (
    id integer NOT NULL,
    family character varying(16)
);


--
-- Name: lexicon_families_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.lexicon_families_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: lexicon_families_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.lexicon_families_id_seq OWNED BY public.lexicon_families.id;


--
-- Name: players; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.players (
    id integer NOT NULL,
    handle character varying(16) NOT NULL,
    woogles character varying(16),
    surname character varying(32),
    first character varying(32),
    middle character varying(32),
    email character varying(64),
    phone bigint,
    lex_family character varying(16)
);


--
-- Name: players_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.players_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: players_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.players_id_seq OWNED BY public.players.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    date date NOT NULL,
    live boolean DEFAULT true,
    location character varying(16),
    notes text
);


--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- Name: words; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.words (
    word character varying(28),
    twl98 boolean DEFAULT false,
    owl1 boolean DEFAULT false,
    owl2 boolean DEFAULT false,
    owl2_1 boolean DEFAULT false,
    twl14 boolean DEFAULT false,
    twl16 boolean DEFAULT false,
    nwl18 boolean DEFAULT false,
    nwl20 boolean DEFAULT false,
    csw07 boolean DEFAULT false,
    csw12 boolean DEFAULT false,
    csw15 boolean DEFAULT false,
    csw19 boolean DEFAULT false,
    csw21 boolean DEFAULT false
);


--
-- Name: games id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games ALTER COLUMN id SET DEFAULT nextval('public.games_id_seq'::regclass);


--
-- Name: lexica id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexica ALTER COLUMN id SET DEFAULT nextval('public.lexica_id_seq'::regclass);


--
-- Name: lexicon_families id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexicon_families ALTER COLUMN id SET DEFAULT nextval('public.lexicon_families_id_seq'::regclass);


--
-- Name: players id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players ALTER COLUMN id SET DEFAULT nextval('public.players_id_seq'::regclass);


--
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (id);


--
-- Name: lexica lexica_lexicon_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexica
    ADD CONSTRAINT lexica_lexicon_key UNIQUE (lexicon);


--
-- Name: lexica lexica_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexica
    ADD CONSTRAINT lexica_pkey PRIMARY KEY (id);


--
-- Name: lexicon_families lexicon_families_family_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexicon_families
    ADD CONSTRAINT lexicon_families_family_key UNIQUE (family);


--
-- Name: lexicon_families lexicon_families_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.lexicon_families
    ADD CONSTRAINT lexicon_families_pkey PRIMARY KEY (id);


--
-- Name: players players_handle_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_handle_key UNIQUE (handle);


--
-- Name: players players_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: words words_word_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.words
    ADD CONSTRAINT words_word_key UNIQUE (word);


--
-- Name: words words_word_key1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.words
    ADD CONSTRAINT words_word_key1 UNIQUE (word);


--
-- Name: bingos bingos_game_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bingos
    ADD CONSTRAINT bingos_game_fkey FOREIGN KEY (game) REFERENCES public.games(id);


--
-- Name: bingos bingos_player_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bingos
    ADD CONSTRAINT bingos_player_fkey FOREIGN KEY (player) REFERENCES public.players(id);


--
-- Name: games games_blank_1_p_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_blank_1_p_fkey FOREIGN KEY (blank_1_p) REFERENCES public.players(id);


--
-- Name: games games_blank_2_p_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_blank_2_p_fkey FOREIGN KEY (blank_2_p) REFERENCES public.players(id);


--
-- Name: games games_player_1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_player_1_fkey FOREIGN KEY (player_1) REFERENCES public.players(id);


--
-- Name: games games_player_2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_player_2_fkey FOREIGN KEY (player_2) REFERENCES public.players(id);


--
-- Name: games games_session_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_session_fkey FOREIGN KEY (session) REFERENCES public.sessions(id);


--
-- PostgreSQL database dump complete
--

