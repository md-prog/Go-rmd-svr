--
-- PostgreSQL database dump
--

-- Dumped from database version 9.4.5
-- Dumped by pg_dump version 9.6.2

-- Started on 2017-03-15 09:32:06 EDT

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12123)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2626 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 174 (class 1259 OID 33546)
-- Name: contact_instance; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE contact_instance (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    contact_type text NOT NULL,
    contact_id integer NOT NULL,
    user_facility_id integer,
    user_contract_id integer
);


ALTER TABLE contact_instance OWNER TO chrislewis;

--
-- TOC entry 173 (class 1259 OID 33544)
-- Name: contact_instance_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE contact_instance_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE contact_instance_id_seq OWNER TO chrislewis;

--
-- TOC entry 2627 (class 0 OID 0)
-- Dependencies: 173
-- Name: contact_instance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE contact_instance_id_seq OWNED BY contact_instance.id;


--
-- TOC entry 176 (class 1259 OID 33558)
-- Name: contract_facilities; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE contract_facilities (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    facilities_id integer NOT NULL,
    contracts_id integer NOT NULL
);


ALTER TABLE contract_facilities OWNER TO chrislewis;

--
-- TOC entry 175 (class 1259 OID 33556)
-- Name: contract_facilities_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE contract_facilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE contract_facilities_id_seq OWNER TO chrislewis;

--
-- TOC entry 2628 (class 0 OID 0)
-- Dependencies: 175
-- Name: contract_facilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE contract_facilities_id_seq OWNED BY contract_facilities.id;


--
-- TOC entry 178 (class 1259 OID 33567)
-- Name: contract_status; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE contract_status (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    lead_no_contact timestamp without time zone,
    initial_contact timestamp without time zone,
    interested timestamp without time zone,
    not_interested_atm timestamp without time zone,
    proposal_sent timestamp without time zone,
    contract_sent timestamp without time zone,
    contract_signed timestamp without time zone,
    working timestamp without time zone,
    dead timestamp without time zone,
    undesirable timestamp without time zone,
    contract_id integer NOT NULL
);


ALTER TABLE contract_status OWNER TO chrislewis;

--
-- TOC entry 177 (class 1259 OID 33565)
-- Name: contract_status_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE contract_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE contract_status_id_seq OWNER TO chrislewis;

--
-- TOC entry 2629 (class 0 OID 0)
-- Dependencies: 177
-- Name: contract_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE contract_status_id_seq OWNED BY contract_status.id;


--
-- TOC entry 180 (class 1259 OID 33576)
-- Name: contracts; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE contracts (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    renewal_date timestamp without time zone,
    lead_source_id integer,
    file text,
    last_contact_id integer NOT NULL
);


ALTER TABLE contracts OWNER TO chrislewis;

--
-- TOC entry 179 (class 1259 OID 33574)
-- Name: contracts_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE contracts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE contracts_id_seq OWNER TO chrislewis;

--
-- TOC entry 2630 (class 0 OID 0)
-- Dependencies: 179
-- Name: contracts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE contracts_id_seq OWNED BY contracts.id;


--
-- TOC entry 182 (class 1259 OID 33588)
-- Name: emrs; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE emrs (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text
);


ALTER TABLE emrs OWNER TO chrislewis;

--
-- TOC entry 181 (class 1259 OID 33586)
-- Name: emrs_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE emrs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE emrs_id_seq OWNER TO chrislewis;

--
-- TOC entry 2631 (class 0 OID 0)
-- Dependencies: 181
-- Name: emrs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE emrs_id_seq OWNED BY emrs.id;


--
-- TOC entry 184 (class 1259 OID 33600)
-- Name: facilities; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE facilities (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text,
    hospitalist_id integer,
    facility_type_id integer NOT NULL,
    emrs_id integer,
    address text,
    city text,
    zip text,
    lat double precision,
    lon double precision,
    facility_stats_id integer NOT NULL,
    primary_contact_id integer NOT NULL,
    secondary_contact_id integer,
    admission_orders boolean DEFAULT false,
    state_id integer,
    last_contact_id integer
);


ALTER TABLE facilities OWNER TO chrislewis;

--
-- TOC entry 183 (class 1259 OID 33598)
-- Name: facilities_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE facilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE facilities_id_seq OWNER TO chrislewis;

--
-- TOC entry 2632 (class 0 OID 0)
-- Dependencies: 183
-- Name: facilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE facilities_id_seq OWNED BY facilities.id;


--
-- TOC entry 185 (class 1259 OID 33611)
-- Name: facility_area_hobbies; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE facility_area_hobbies (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    facility_id integer NOT NULL,
    hobby_id integer NOT NULL
);


ALTER TABLE facility_area_hobbies OWNER TO chrislewis;

--
-- TOC entry 187 (class 1259 OID 33619)
-- Name: facility_stats; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE facility_stats (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    entry_date timestamp without time zone DEFAULT now() NOT NULL,
    bed_count_ft integer DEFAULT 0,
    bed_count_main integer DEFAULT 0,
    bed_count_overflow integer,
    annual_volume integer DEFAULT 0,
    provider_in_triage boolean DEFAULT false,
    hours_physician integer DEFAULT 0,
    hours_apc integer DEFAULT 0,
    shifts_physician integer DEFAULT 0,
    shifts_apc integer DEFAULT 0,
    annual_hospitalizations integer DEFAULT 0,
    annual_obs integer DEFAULT 0,
    annual_admit integer DEFAULT 0,
    lwbs integer DEFAULT 0,
    lwot integer DEFAULT 0,
    ft_hours_start timestamp without time zone,
    ft_hours_end timestamp without time zone
);


ALTER TABLE facility_stats OWNER TO chrislewis;

--
-- TOC entry 186 (class 1259 OID 33617)
-- Name: facility_stats_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE facility_stats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE facility_stats_id_seq OWNER TO chrislewis;

--
-- TOC entry 2633 (class 0 OID 0)
-- Dependencies: 186
-- Name: facility_stats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE facility_stats_id_seq OWNED BY facility_stats.id;


--
-- TOC entry 189 (class 1259 OID 33642)
-- Name: facility_types; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE facility_types (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text NOT NULL
);


ALTER TABLE facility_types OWNER TO chrislewis;

--
-- TOC entry 188 (class 1259 OID 33640)
-- Name: facility_types_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE facility_types_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE facility_types_id_seq OWNER TO chrislewis;

--
-- TOC entry 2634 (class 0 OID 0)
-- Dependencies: 188
-- Name: facility_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE facility_types_id_seq OWNED BY facility_types.id;


--
-- TOC entry 191 (class 1259 OID 33654)
-- Name: hobbies; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE hobbies (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text NOT NULL
);


ALTER TABLE hobbies OWNER TO chrislewis;

--
-- TOC entry 190 (class 1259 OID 33652)
-- Name: hobbies_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE hobbies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE hobbies_id_seq OWNER TO chrislewis;

--
-- TOC entry 2635 (class 0 OID 0)
-- Dependencies: 190
-- Name: hobbies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE hobbies_id_seq OWNED BY hobbies.id;


--
-- TOC entry 193 (class 1259 OID 33666)
-- Name: hospitalists; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE hospitalists (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text,
    primary_contact_id integer
);


ALTER TABLE hospitalists OWNER TO chrislewis;

--
-- TOC entry 192 (class 1259 OID 33664)
-- Name: hospitalists_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE hospitalists_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE hospitalists_id_seq OWNER TO chrislewis;

--
-- TOC entry 2636 (class 0 OID 0)
-- Dependencies: 192
-- Name: hospitalists_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE hospitalists_id_seq OWNED BY hospitalists.id;


--
-- TOC entry 195 (class 1259 OID 33678)
-- Name: people; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE people (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    middle_name text,
    title text,
    email text NOT NULL,
    cell_phone text,
    home_phone text,
    office_phone text,
    fax text,
    address text,
    zip text,
    city text,
    lat double precision,
    lon double precision,
    state_id integer,
    recruiter_id integer,
    provider_id integer
);


ALTER TABLE people OWNER TO chrislewis;

--
-- TOC entry 194 (class 1259 OID 33676)
-- Name: people_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE people_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE people_id_seq OWNER TO chrislewis;

--
-- TOC entry 2637 (class 0 OID 0)
-- Dependencies: 194
-- Name: people_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE people_id_seq OWNED BY people.id;


--
-- TOC entry 197 (class 1259 OID 33690)
-- Name: provider_certifications; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_certifications (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    expiration_date timestamp without time zone NOT NULL,
    file text,
    provider_id integer NOT NULL
);


ALTER TABLE provider_certifications OWNER TO chrislewis;

--
-- TOC entry 196 (class 1259 OID 33688)
-- Name: provider_certifications_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_certifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_certifications_id_seq OWNER TO chrislewis;

--
-- TOC entry 2638 (class 0 OID 0)
-- Dependencies: 196
-- Name: provider_certifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_certifications_id_seq OWNED BY provider_certifications.id;


--
-- TOC entry 199 (class 1259 OID 33702)
-- Name: provider_contracts; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_contracts (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    start_date timestamp without time zone,
    renewal_date timestamp without time zone,
    file text,
    provider_id integer NOT NULL,
    user_facilities_id integer,
    user_contracts_id integer
);


ALTER TABLE provider_contracts OWNER TO chrislewis;

--
-- TOC entry 198 (class 1259 OID 33700)
-- Name: provider_contracts_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_contracts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_contracts_id_seq OWNER TO chrislewis;

--
-- TOC entry 2639 (class 0 OID 0)
-- Dependencies: 198
-- Name: provider_contracts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_contracts_id_seq OWNED BY provider_contracts.id;


--
-- TOC entry 201 (class 1259 OID 33714)
-- Name: provider_facilities; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_facilities (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    provider_id integer NOT NULL,
    facility_id integer NOT NULL,
    credentialed_date timestamp without time zone,
    renewal_date timestamp without time zone
);


ALTER TABLE provider_facilities OWNER TO chrislewis;

--
-- TOC entry 200 (class 1259 OID 33712)
-- Name: provider_facilities_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_facilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_facilities_id_seq OWNER TO chrislewis;

--
-- TOC entry 2640 (class 0 OID 0)
-- Dependencies: 200
-- Name: provider_facilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_facilities_id_seq OWNED BY provider_facilities.id;


--
-- TOC entry 203 (class 1259 OID 33723)
-- Name: provider_hobbies; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_hobbies (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    provider_id integer NOT NULL,
    hobby_id integer NOT NULL
);


ALTER TABLE provider_hobbies OWNER TO chrislewis;

--
-- TOC entry 202 (class 1259 OID 33721)
-- Name: provider_hobbies_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_hobbies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_hobbies_id_seq OWNER TO chrislewis;

--
-- TOC entry 2641 (class 0 OID 0)
-- Dependencies: 202
-- Name: provider_hobbies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_hobbies_id_seq OWNED BY provider_hobbies.id;


--
-- TOC entry 205 (class 1259 OID 33732)
-- Name: provider_licenses; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_licenses (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    state_id integer NOT NULL,
    license_number integer NOT NULL,
    file text,
    provider_id integer NOT NULL
);


ALTER TABLE provider_licenses OWNER TO chrislewis;

--
-- TOC entry 204 (class 1259 OID 33730)
-- Name: provider_licenses_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_licenses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_licenses_id_seq OWNER TO chrislewis;

--
-- TOC entry 2642 (class 0 OID 0)
-- Dependencies: 204
-- Name: provider_licenses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_licenses_id_seq OWNED BY provider_licenses.id;


--
-- TOC entry 207 (class 1259 OID 33744)
-- Name: provider_status; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE provider_status (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    lead_no_contact timestamp without time zone,
    initial_contact timestamp without time zone,
    interested timestamp without time zone,
    not_interested_atm timestamp without time zone,
    verbal_commitment timestamp without time zone,
    contract_sent timestamp without time zone,
    contract_signed timestamp without time zone,
    dead timestamp without time zone,
    do_not_use timestamp without time zone,
    provider_id integer NOT NULL
);


ALTER TABLE provider_status OWNER TO chrislewis;

--
-- TOC entry 206 (class 1259 OID 33742)
-- Name: provider_status_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE provider_status_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE provider_status_id_seq OWNER TO chrislewis;

--
-- TOC entry 2643 (class 0 OID 0)
-- Dependencies: 206
-- Name: provider_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE provider_status_id_seq OWNED BY provider_status.id;


--
-- TOC entry 209 (class 1259 OID 33753)
-- Name: providers; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE providers (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    npi integer,
    w9 text,
    direct_deposit_form text,
    lead_source_id integer,
    responsible_recruiter_id integer,
    hourly_rate integer,
    desired_shifts_month integer,
    max_shifts_month integer,
    min_shifts_month integer,
    full_time boolean,
    part_time boolean,
    prn boolean,
    retired boolean,
    notes text,
    insurance_certificate text,
    tb_expiration timestamp without time zone,
    tb_file text,
    flu_expiration timestamp without time zone,
    flu_file text
);


ALTER TABLE providers OWNER TO chrislewis;

--
-- TOC entry 208 (class 1259 OID 33751)
-- Name: providers_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE providers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE providers_id_seq OWNER TO chrislewis;

--
-- TOC entry 2644 (class 0 OID 0)
-- Dependencies: 208
-- Name: providers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE providers_id_seq OWNED BY providers.id;


--
-- TOC entry 211 (class 1259 OID 33765)
-- Name: recruiters; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE recruiters (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE recruiters OWNER TO chrislewis;

--
-- TOC entry 210 (class 1259 OID 33763)
-- Name: recruiters_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE recruiters_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE recruiters_id_seq OWNER TO chrislewis;

--
-- TOC entry 2645 (class 0 OID 0)
-- Dependencies: 210
-- Name: recruiters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE recruiters_id_seq OWNED BY recruiters.id;


--
-- TOC entry 213 (class 1259 OID 33774)
-- Name: states; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE states (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    abbreviation text NOT NULL
);


ALTER TABLE states OWNER TO chrislewis;

--
-- TOC entry 212 (class 1259 OID 33772)
-- Name: states_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE states_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE states_id_seq OWNER TO chrislewis;

--
-- TOC entry 2646 (class 0 OID 0)
-- Dependencies: 212
-- Name: states_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE states_id_seq OWNED BY states.id;


--
-- TOC entry 215 (class 1259 OID 33786)
-- Name: user_contacts; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE user_contacts (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL,
    people_id integer NOT NULL
);


ALTER TABLE user_contacts OWNER TO chrislewis;

--
-- TOC entry 214 (class 1259 OID 33784)
-- Name: user_contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE user_contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_contacts_id_seq OWNER TO chrislewis;

--
-- TOC entry 2647 (class 0 OID 0)
-- Dependencies: 214
-- Name: user_contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE user_contacts_id_seq OWNED BY user_contacts.id;


--
-- TOC entry 217 (class 1259 OID 33795)
-- Name: user_contracts; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE user_contracts (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL,
    contracts_id integer NOT NULL
);


ALTER TABLE user_contracts OWNER TO chrislewis;

--
-- TOC entry 216 (class 1259 OID 33793)
-- Name: user_contracts_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE user_contracts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_contracts_id_seq OWNER TO chrislewis;

--
-- TOC entry 2648 (class 0 OID 0)
-- Dependencies: 216
-- Name: user_contracts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE user_contracts_id_seq OWNED BY user_contracts.id;


--
-- TOC entry 219 (class 1259 OID 33804)
-- Name: user_facilities; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE user_facilities (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL,
    facility_id integer NOT NULL
);


ALTER TABLE user_facilities OWNER TO chrislewis;

--
-- TOC entry 218 (class 1259 OID 33802)
-- Name: user_facilities_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE user_facilities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE user_facilities_id_seq OWNER TO chrislewis;

--
-- TOC entry 2649 (class 0 OID 0)
-- Dependencies: 218
-- Name: user_facilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE user_facilities_id_seq OWNED BY user_facilities.id;


--
-- TOC entry 221 (class 1259 OID 33813)
-- Name: users; Type: TABLE; Schema: public; Owner: chrislewis
--

CREATE TABLE users (
    id integer NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    people_id integer,
    email text NOT NULL,
    password text NOT NULL
);


ALTER TABLE users OWNER TO chrislewis;

--
-- TOC entry 220 (class 1259 OID 33811)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: chrislewis
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO chrislewis;

--
-- TOC entry 2650 (class 0 OID 0)
-- Dependencies: 220
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: chrislewis
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- TOC entry 2303 (class 2604 OID 33549)
-- Name: contact_instance id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contact_instance ALTER COLUMN id SET DEFAULT nextval('contact_instance_id_seq'::regclass);


--
-- TOC entry 2305 (class 2604 OID 33561)
-- Name: contract_facilities id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_facilities ALTER COLUMN id SET DEFAULT nextval('contract_facilities_id_seq'::regclass);


--
-- TOC entry 2307 (class 2604 OID 33570)
-- Name: contract_status id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_status ALTER COLUMN id SET DEFAULT nextval('contract_status_id_seq'::regclass);


--
-- TOC entry 2309 (class 2604 OID 33579)
-- Name: contracts id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contracts ALTER COLUMN id SET DEFAULT nextval('contracts_id_seq'::regclass);


--
-- TOC entry 2311 (class 2604 OID 33591)
-- Name: emrs id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY emrs ALTER COLUMN id SET DEFAULT nextval('emrs_id_seq'::regclass);


--
-- TOC entry 2313 (class 2604 OID 33603)
-- Name: facilities id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities ALTER COLUMN id SET DEFAULT nextval('facilities_id_seq'::regclass);


--
-- TOC entry 2317 (class 2604 OID 33622)
-- Name: facility_stats id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_stats ALTER COLUMN id SET DEFAULT nextval('facility_stats_id_seq'::regclass);


--
-- TOC entry 2333 (class 2604 OID 33645)
-- Name: facility_types id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_types ALTER COLUMN id SET DEFAULT nextval('facility_types_id_seq'::regclass);


--
-- TOC entry 2335 (class 2604 OID 33657)
-- Name: hobbies id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hobbies ALTER COLUMN id SET DEFAULT nextval('hobbies_id_seq'::regclass);


--
-- TOC entry 2337 (class 2604 OID 33669)
-- Name: hospitalists id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hospitalists ALTER COLUMN id SET DEFAULT nextval('hospitalists_id_seq'::regclass);


--
-- TOC entry 2339 (class 2604 OID 33681)
-- Name: people id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY people ALTER COLUMN id SET DEFAULT nextval('people_id_seq'::regclass);


--
-- TOC entry 2341 (class 2604 OID 33693)
-- Name: provider_certifications id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_certifications ALTER COLUMN id SET DEFAULT nextval('provider_certifications_id_seq'::regclass);


--
-- TOC entry 2343 (class 2604 OID 33705)
-- Name: provider_contracts id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_contracts ALTER COLUMN id SET DEFAULT nextval('provider_contracts_id_seq'::regclass);


--
-- TOC entry 2345 (class 2604 OID 33717)
-- Name: provider_facilities id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_facilities ALTER COLUMN id SET DEFAULT nextval('provider_facilities_id_seq'::regclass);


--
-- TOC entry 2347 (class 2604 OID 33726)
-- Name: provider_hobbies id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_hobbies ALTER COLUMN id SET DEFAULT nextval('provider_hobbies_id_seq'::regclass);


--
-- TOC entry 2349 (class 2604 OID 33735)
-- Name: provider_licenses id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_licenses ALTER COLUMN id SET DEFAULT nextval('provider_licenses_id_seq'::regclass);


--
-- TOC entry 2351 (class 2604 OID 33747)
-- Name: provider_status id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_status ALTER COLUMN id SET DEFAULT nextval('provider_status_id_seq'::regclass);


--
-- TOC entry 2353 (class 2604 OID 33756)
-- Name: providers id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY providers ALTER COLUMN id SET DEFAULT nextval('providers_id_seq'::regclass);


--
-- TOC entry 2355 (class 2604 OID 33768)
-- Name: recruiters id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY recruiters ALTER COLUMN id SET DEFAULT nextval('recruiters_id_seq'::regclass);


--
-- TOC entry 2357 (class 2604 OID 33777)
-- Name: states id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY states ALTER COLUMN id SET DEFAULT nextval('states_id_seq'::regclass);


--
-- TOC entry 2359 (class 2604 OID 33789)
-- Name: user_contacts id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contacts ALTER COLUMN id SET DEFAULT nextval('user_contacts_id_seq'::regclass);


--
-- TOC entry 2361 (class 2604 OID 33798)
-- Name: user_contracts id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contracts ALTER COLUMN id SET DEFAULT nextval('user_contracts_id_seq'::regclass);


--
-- TOC entry 2363 (class 2604 OID 33807)
-- Name: user_facilities id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_facilities ALTER COLUMN id SET DEFAULT nextval('user_facilities_id_seq'::regclass);


--
-- TOC entry 2365 (class 2604 OID 33816)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- TOC entry 2571 (class 0 OID 33546)
-- Dependencies: 174
-- Data for Name: contact_instance; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY contact_instance (id, updated_at, contact_type, contact_id, user_facility_id, user_contract_id) FROM stdin;
1	2017-03-14 22:38:00.485862	email	2	\N	\N
2	2017-03-14 22:38:00.485862	phone	2	\N	\N
3	2017-03-14 22:38:00.485862	phone	6	\N	\N
4	2017-03-14 22:38:00.485862	in_person	9	\N	\N
\.


--
-- TOC entry 2651 (class 0 OID 0)
-- Dependencies: 173
-- Name: contact_instance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('contact_instance_id_seq', 4, true);


--
-- TOC entry 2573 (class 0 OID 33558)
-- Dependencies: 176
-- Data for Name: contract_facilities; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY contract_facilities (id, updated_at, facilities_id, contracts_id) FROM stdin;
1	2017-03-14 23:27:48.431102	1	1
2	2017-03-14 23:27:48.431102	4	2
3	2017-03-14 23:27:48.431102	2	1
4	2017-03-14 23:27:48.431102	2	3
5	2017-03-14 23:27:48.431102	8	4
\.


--
-- TOC entry 2652 (class 0 OID 0)
-- Dependencies: 175
-- Name: contract_facilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('contract_facilities_id_seq', 5, true);


--
-- TOC entry 2575 (class 0 OID 33567)
-- Dependencies: 178
-- Data for Name: contract_status; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY contract_status (id, updated_at, lead_no_contact, initial_contact, interested, not_interested_atm, proposal_sent, contract_sent, contract_signed, working, dead, undesirable, contract_id) FROM stdin;
1	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	1
2	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	2
3	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	3
4	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	4
5	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	5
6	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	6
7	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	7
8	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	8
9	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	9
10	2017-03-14 23:22:22.423245	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
\.


--
-- TOC entry 2653 (class 0 OID 0)
-- Dependencies: 177
-- Name: contract_status_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('contract_status_id_seq', 10, true);


--
-- TOC entry 2577 (class 0 OID 33576)
-- Dependencies: 180
-- Data for Name: contracts; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY contracts (id, updated_at, renewal_date, lead_source_id, file, last_contact_id) FROM stdin;
1	2017-03-07 06:22:27	2016-11-10 01:38:01	5	\N	3
2	2016-10-26 06:28:02	2016-08-28 02:57:42	5	\N	3
3	2016-04-29 08:03:25	2016-07-15 06:01:34	2	\N	3
4	2017-01-11 19:31:01	2017-03-08 04:26:07	1	\N	3
5	2016-08-08 20:05:04	2017-01-21 10:00:39	1	\N	2
6	2016-11-04 23:15:35	2016-10-14 03:33:01	1	\N	3
7	2016-07-04 04:05:59	2016-05-08 03:11:43	2	\N	3
8	2016-05-01 16:24:41	2016-06-04 09:37:44	5	\N	3
9	2016-05-01 15:27:16	2016-04-06 01:56:32	3	\N	3
10	2016-03-31 01:39:54	2016-09-15 02:11:53	2	\N	3
\.


--
-- TOC entry 2654 (class 0 OID 0)
-- Dependencies: 179
-- Name: contracts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('contracts_id_seq', 1, false);


--
-- TOC entry 2579 (class 0 OID 33588)
-- Dependencies: 182
-- Data for Name: emrs; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY emrs (id, updated_at, name) FROM stdin;
1	2016-08-02 18:51:14	Yadel
2	2016-10-09 03:21:52	Browsetype
3	2016-06-29 15:04:41	Zava
4	2016-07-09 19:17:31	Aimbo
5	2016-05-16 00:17:21	Realcube
6	2016-03-25 06:14:24	Roombo
7	2016-03-30 04:56:48	Brainsphere
8	2016-03-15 21:14:07	Buzzster
9	2016-08-08 13:34:19	Avamm
10	2016-03-17 20:13:15	Gevee
\.


--
-- TOC entry 2655 (class 0 OID 0)
-- Dependencies: 181
-- Name: emrs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('emrs_id_seq', 1, false);


--
-- TOC entry 2581 (class 0 OID 33600)
-- Dependencies: 184
-- Data for Name: facilities; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY facilities (id, updated_at, name, hospitalist_id, facility_type_id, emrs_id, address, city, zip, lat, lon, facility_stats_id, primary_contact_id, secondary_contact_id, admission_orders, state_id, last_contact_id) FROM stdin;
1	2016-05-22 23:27:17	Divanoodle	\N	1	\N	76 Bluestem Circle	Taungoo	\N	\N	\N	1	2	\N	t	7	\N
2	2016-08-16 13:00:08	Plajo	\N	2	\N	9 Artisan Drive	Mata de São João	48280-000	\N	\N	10	1	\N	f	6	\N
3	2016-03-14 18:46:53	Youspan	\N	2	\N	08947 Havey Hill	Ishimbay	453209	\N	\N	7	7	\N	f	5	\N
4	2017-02-18 17:33:57	Yamia	\N	2	\N	61 Maryland Junction	Praszka	46-320	\N	\N	2	3	\N	f	10	\N
5	2016-10-20 07:51:11	Ozu	\N	2	\N	764 Mcbride Place	Ciheras	\N	\N	\N	3	2	\N	f	8	\N
6	2016-12-14 08:18:09	Jabbersphere	\N	2	\N	254 Spohn Hill	Dhībīn	\N	\N	\N	7	4	\N	t	2	\N
7	2016-12-05 06:32:29	Feedfire	\N	2	\N	937 Mcguire Alley	Fuzhiping	\N	\N	\N	1	3	\N	f	8	\N
8	2017-01-10 07:19:25	Dazzlesphere	\N	2	\N	66808 Maple Trail	Macau	59500-000	\N	\N	7	6	\N	t	6	\N
9	2017-02-21 23:53:14	Rhyloo	\N	2	\N	7 Summer Ridge Junction	Avignon	84033 CEDEX 3	\N	\N	1	5	\N	f	6	\N
10	2016-04-06 00:37:59	Fivebridge	\N	3	\N	68107 Birchwood Court	Toledo	43605	\N	\N	5	2	\N	t	10	\N
\.


--
-- TOC entry 2656 (class 0 OID 0)
-- Dependencies: 183
-- Name: facilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('facilities_id_seq', 1, false);


--
-- TOC entry 2582 (class 0 OID 33611)
-- Dependencies: 185
-- Data for Name: facility_area_hobbies; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY facility_area_hobbies (id, updated_at, facility_id, hobby_id) FROM stdin;
\.


--
-- TOC entry 2584 (class 0 OID 33619)
-- Dependencies: 187
-- Data for Name: facility_stats; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY facility_stats (id, updated_at, entry_date, bed_count_ft, bed_count_main, bed_count_overflow, annual_volume, provider_in_triage, hours_physician, hours_apc, shifts_physician, shifts_apc, annual_hospitalizations, annual_obs, annual_admit, lwbs, lwot, ft_hours_start, ft_hours_end) FROM stdin;
1	2017-03-14 22:16:31.940441	2017-03-14 22:16:31.940441	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
2	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
3	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
4	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
5	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
6	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
7	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
8	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
9	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
10	2017-03-14 22:16:44.540922	2017-03-14 22:16:44.540922	0	0	\N	0	f	0	0	0	0	0	0	0	0	0	\N	\N
\.


--
-- TOC entry 2657 (class 0 OID 0)
-- Dependencies: 186
-- Name: facility_stats_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('facility_stats_id_seq', 10, true);


--
-- TOC entry 2586 (class 0 OID 33642)
-- Dependencies: 189
-- Data for Name: facility_types; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY facility_types (id, updated_at, name) FROM stdin;
1	2017-01-30 13:56:13	Urgent Care
2	2016-08-15 00:06:31	ER
3	2016-08-15 00:06:31	Free Standing ER
\.


--
-- TOC entry 2658 (class 0 OID 0)
-- Dependencies: 188
-- Name: facility_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('facility_types_id_seq', 1, true);


--
-- TOC entry 2588 (class 0 OID 33654)
-- Dependencies: 191
-- Data for Name: hobbies; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY hobbies (id, updated_at, name) FROM stdin;
1	2016-08-31 09:51:41	piano
3	2016-12-15 05:53:07	crocheting
4	2017-01-01 08:45:49	golf
\.


--
-- TOC entry 2659 (class 0 OID 0)
-- Dependencies: 190
-- Name: hobbies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('hobbies_id_seq', 1, false);


--
-- TOC entry 2590 (class 0 OID 33666)
-- Dependencies: 193
-- Data for Name: hospitalists; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY hospitalists (id, updated_at, name, primary_contact_id) FROM stdin;
1	2016-09-04 22:40:47	Shuffletag	10
2	2017-02-16 07:08:44	Wikizz	5
3	2016-04-14 07:48:52	Quinu	10
4	2016-10-04 21:05:28	Trupe	9
5	2016-06-12 03:06:46	Cogilith	10
6	2016-03-22 08:06:14	Oyoloo	6
7	2016-05-29 20:02:50	Devpulse	5
8	2016-03-28 23:26:50	Devcast	1
9	2016-04-16 07:15:44	Skilith	9
10	2016-10-08 12:39:18	Quaxo	1
\.


--
-- TOC entry 2660 (class 0 OID 0)
-- Dependencies: 192
-- Name: hospitalists_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('hospitalists_id_seq', 1, false);


--
-- TOC entry 2592 (class 0 OID 33678)
-- Dependencies: 195
-- Data for Name: people; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY people (id, updated_at, first_name, last_name, middle_name, title, email, cell_phone, home_phone, office_phone, fax, address, zip, city, lat, lon, state_id, recruiter_id, provider_id) FROM stdin;
1	2016-04-07 15:52:07	Steve	Elliott	Donald	Ms	delliott0@cdbaby.com	58-(538)946-7054	62-(122)120-6164	502-(449)323-5597	92-(940)601-7212	4735 Shoshone Junction	\N	San José de Barlovento	\N	\N	8	\N	\N
2	2016-09-03 00:24:28	Peter	Hamilton	Karen	Ms	khamilton1@people.com.cn	86-(841)590-1742	86-(984)262-4251	55-(579)949-5032	351-(398)194-2187	8 Almo Place	\N	Chifeng	\N	\N	1	\N	\N
3	2016-04-15 20:28:49	Ann	Fields	Tina	Mrs	tfields2@guardian.co.uk	7-(949)767-2980	49-(587)710-4364	86-(315)428-8600	86-(867)291-1539	6 Bunting Drive	624082	Sysert’	\N	\N	7	\N	\N
4	2016-05-28 05:13:37	Julia	Hughes	Henry	Mr	hhughes3@goo.gl	62-(559)236-9867	62-(839)561-2088	62-(282)250-4452	62-(663)694-9388	884 Coolidge Junction	\N	Jalgung	\N	\N	5	\N	\N
5	2017-01-24 18:14:34	Nancy	Burton	Robin	Mr	rburton4@washington.edu	56-(232)518-6012	420-(647)113-4405	1-(304)896-8142	86-(324)167-7483	2833 Talisman Crossing	\N	Coihaique	\N	\N	2	\N	\N
6	2017-01-26 00:57:02	Gloria	Black	Kevin	Dr	kblack5@earthlink.net	383-(739)561-3972	86-(557)664-1358	62-(439)320-9347	355-(255)723-1556	1 Elmside Point	\N	Klokot	\N	\N	8	\N	\N
7	2017-03-12 09:05:26	Donald	Schmidt	Emily	Ms	eschmidt6@fotki.com	63-(824)470-9193	371-(759)983-3529	7-(802)460-2393	86-(490)169-9513	98078 Johnson Hill	6806	Masaguisi	\N	\N	2	\N	\N
8	2016-11-20 04:04:58	Mary	Johnson	Phyllis	Mr	pjohnson7@arizona.edu	46-(289)560-9864	380-(347)275-8284	49-(762)917-1206	236-(855)962-9103	455 Forest Dale Alley	355 92	Växjö	\N	\N	2	\N	\N
9	2016-12-17 00:21:49	Diana	Crawford	James	Honorable	jcrawford8@eepurl.com	355-(634)261-8518	51-(112)324-4695	81-(981)653-2713	66-(174)148-4793	3 Lake View Terrace	\N	Hoçisht	\N	\N	2	\N	\N
10	2016-11-07 22:57:15	Martha	Foster	Susan	Mrs	sfoster9@ovh.net	7-(581)496-9663	47-(126)260-8937	355-(797)566-0161	45-(727)827-8043	61 Sunnyside Point	442064	Sosnovka	\N	\N	9	\N	\N
\.


--
-- TOC entry 2661 (class 0 OID 0)
-- Dependencies: 194
-- Name: people_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('people_id_seq', 1, false);


--
-- TOC entry 2594 (class 0 OID 33690)
-- Dependencies: 197
-- Data for Name: provider_certifications; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_certifications (id, updated_at, name, expiration_date, file, provider_id) FROM stdin;
1	2016-03-14 05:32:41	NRP	2016-12-30 15:30:07	\N	9
2	2016-09-09 12:07:38	NRP	2017-01-28 13:43:58	\N	8
3	2016-07-20 04:43:58	NRP	2016-03-15 01:36:39	\N	10
4	2016-08-26 09:55:44	ACLS	2016-08-18 03:33:51	\N	4
5	2016-08-15 21:37:01	NRP	2016-06-09 01:18:03	\N	5
6	2016-04-28 19:21:36	ACLS	2016-12-01 13:30:01	\N	9
7	2017-01-22 08:52:11	ACLS	2016-07-17 07:19:59	\N	3
8	2016-12-30 04:02:41	NRP	2016-05-25 10:53:32	\N	4
9	2016-09-10 09:18:32	NRP	2017-01-17 03:20:10	\N	5
10	2016-03-28 02:55:53	PALS	2016-04-08 13:36:42	\N	2
\.


--
-- TOC entry 2662 (class 0 OID 0)
-- Dependencies: 196
-- Name: provider_certifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_certifications_id_seq', 1, false);


--
-- TOC entry 2596 (class 0 OID 33702)
-- Dependencies: 199
-- Data for Name: provider_contracts; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_contracts (id, updated_at, start_date, renewal_date, file, provider_id, user_facilities_id, user_contracts_id) FROM stdin;
\.


--
-- TOC entry 2663 (class 0 OID 0)
-- Dependencies: 198
-- Name: provider_contracts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_contracts_id_seq', 1, false);


--
-- TOC entry 2598 (class 0 OID 33714)
-- Dependencies: 201
-- Data for Name: provider_facilities; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_facilities (id, updated_at, provider_id, facility_id, credentialed_date, renewal_date) FROM stdin;
\.


--
-- TOC entry 2664 (class 0 OID 0)
-- Dependencies: 200
-- Name: provider_facilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_facilities_id_seq', 1, false);


--
-- TOC entry 2600 (class 0 OID 33723)
-- Dependencies: 203
-- Data for Name: provider_hobbies; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_hobbies (id, updated_at, provider_id, hobby_id) FROM stdin;
\.


--
-- TOC entry 2665 (class 0 OID 0)
-- Dependencies: 202
-- Name: provider_hobbies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_hobbies_id_seq', 1, false);


--
-- TOC entry 2602 (class 0 OID 33732)
-- Dependencies: 205
-- Data for Name: provider_licenses; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_licenses (id, updated_at, state_id, license_number, file, provider_id) FROM stdin;
1	2017-03-12 18:06:46	10	46130643	\N	5
2	2016-11-09 02:47:05	29	85016238	\N	7
3	2016-06-10 15:21:46	26	84627926	\N	6
4	2016-06-09 04:15:18	13	43386851	\N	7
5	2017-02-21 18:56:43	6	81262958	\N	5
6	2016-10-14 13:52:24	20	96740303	\N	3
7	2016-12-21 10:44:27	6	89327042	\N	6
8	2017-01-03 14:15:15	1	23781641	\N	2
9	2016-10-28 03:50:32	21	14262060	\N	7
10	2016-12-16 05:45:19	9	39452053	\N	6
11	2016-04-08 10:08:19	5	23035953	\N	9
12	2017-01-20 13:45:19	8	95997913	\N	5
13	2016-04-25 12:30:24	25	53126903	\N	9
14	2016-12-15 04:08:33	9	19852492	\N	9
15	2017-01-23 07:21:56	23	55554130	\N	9
16	2016-03-15 11:05:23	1	46548384	\N	2
17	2017-01-26 11:06:17	25	78841158	\N	5
18	2016-12-27 05:03:49	25	19180281	\N	7
19	2016-08-14 05:51:58	7	75004625	\N	9
20	2016-09-03 23:54:41	30	33058484	\N	3
\.


--
-- TOC entry 2666 (class 0 OID 0)
-- Dependencies: 204
-- Name: provider_licenses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_licenses_id_seq', 1, false);


--
-- TOC entry 2604 (class 0 OID 33744)
-- Dependencies: 207
-- Data for Name: provider_status; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY provider_status (id, updated_at, lead_no_contact, initial_contact, interested, not_interested_atm, verbal_commitment, contract_sent, contract_signed, dead, do_not_use, provider_id) FROM stdin;
1	2016-06-24 15:28:09	\N	\N	\N	\N	\N	\N	\N	\N	\N	9
2	2017-02-01 04:39:26	\N	\N	\N	\N	\N	\N	\N	\N	\N	2
3	2016-03-31 08:27:35	\N	\N	\N	\N	\N	\N	\N	\N	\N	7
4	2016-04-26 02:43:49	\N	\N	\N	\N	\N	\N	\N	\N	\N	6
5	2017-03-12 22:33:32	\N	\N	\N	\N	\N	\N	\N	\N	\N	1
6	2016-12-14 20:53:52	\N	\N	\N	\N	\N	\N	\N	\N	\N	6
7	2016-06-16 03:38:14	\N	\N	\N	\N	\N	\N	\N	\N	\N	8
8	2017-01-04 01:46:50	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
9	2017-02-13 05:17:14	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
10	2017-01-08 00:34:40	\N	\N	\N	\N	\N	\N	\N	\N	\N	8
11	2016-06-21 04:47:25	\N	\N	\N	\N	\N	\N	\N	\N	\N	1
12	2016-09-30 00:57:36	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
13	2017-02-25 00:42:26	\N	\N	\N	\N	\N	\N	\N	\N	\N	4
14	2016-04-06 00:49:06	\N	\N	\N	\N	\N	\N	\N	\N	\N	7
15	2016-07-08 11:29:10	\N	\N	\N	\N	\N	\N	\N	\N	\N	1
16	2016-09-21 22:42:57	\N	\N	\N	\N	\N	\N	\N	\N	\N	6
17	2016-09-21 14:37:44	\N	\N	\N	\N	\N	\N	\N	\N	\N	5
18	2016-05-02 08:48:30	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
19	2017-01-07 23:04:41	\N	\N	\N	\N	\N	\N	\N	\N	\N	6
20	2017-02-19 21:19:59	\N	\N	\N	\N	\N	\N	\N	\N	\N	10
\.


--
-- TOC entry 2667 (class 0 OID 0)
-- Dependencies: 206
-- Name: provider_status_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('provider_status_id_seq', 1, false);


--
-- TOC entry 2606 (class 0 OID 33753)
-- Dependencies: 209
-- Data for Name: providers; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY providers (id, updated_at, npi, w9, direct_deposit_form, lead_source_id, responsible_recruiter_id, hourly_rate, desired_shifts_month, max_shifts_month, min_shifts_month, full_time, part_time, prn, retired, notes, insurance_certificate, tb_expiration, tb_file, flu_expiration, flu_file) FROM stdin;
2	2016-03-25 18:57:22	7	\N	\N	10	\N	9	6	6	5	f	f	t	f	\N	\N	2017-02-21 22:34:26	\N	2016-11-29 09:58:53	\N
3	2016-07-14 04:39:27	8	\N	\N	8	\N	7	1	8	2	f	t	t	f	\N	\N	2017-02-12 00:51:55	\N	2016-12-19 14:00:25	\N
4	2016-05-09 14:24:19	6	\N	\N	8	\N	5	1	5	5	f	t	t	f	\N	\N	2017-02-12 06:42:42	\N	2017-02-13 08:53:14	\N
5	2016-04-26 14:55:30	10	\N	\N	5	\N	5	1	7	8	t	t	f	f	\N	\N	2016-10-08 19:00:44	\N	2016-08-21 19:56:38	\N
6	2016-07-30 17:32:50	5	\N	\N	8	\N	1	10	2	3	t	f	t	f	\N	\N	2016-10-23 13:05:21	\N	2016-10-03 15:56:30	\N
7	2016-05-18 01:34:25	10	\N	\N	7	\N	2	3	6	7	f	t	t	t	\N	\N	2016-10-18 00:31:54	\N	2016-09-21 16:06:36	\N
8	2017-02-25 05:36:53	2	\N	\N	2	\N	9	8	8	1	t	t	t	f	\N	\N	2016-11-22 00:23:15	\N	2017-01-23 04:23:58	\N
9	2017-02-15 01:55:29	2	\N	\N	2	\N	6	3	4	6	t	f	t	f	\N	\N	2017-01-26 15:12:44	\N	2016-08-23 09:04:03	\N
10	2017-03-13 20:18:10	9	\N	\N	4	\N	2	5	5	7	t	f	t	t	\N	\N	2017-02-19 00:33:02	\N	2016-08-24 20:06:21	\N
1	2016-08-25 14:49:29	2	\N	\N	2	\N	6	3	1	8	t	f	t	f	\N	\N	2016-10-03 06:17:04	\N	2016-08-24 13:27:43	\N
\.


--
-- TOC entry 2668 (class 0 OID 0)
-- Dependencies: 208
-- Name: providers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('providers_id_seq', 1, false);


--
-- TOC entry 2608 (class 0 OID 33765)
-- Dependencies: 211
-- Data for Name: recruiters; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY recruiters (id, updated_at) FROM stdin;
\.


--
-- TOC entry 2669 (class 0 OID 0)
-- Dependencies: 210
-- Name: recruiters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('recruiters_id_seq', 1, false);


--
-- TOC entry 2610 (class 0 OID 33774)
-- Dependencies: 213
-- Data for Name: states; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY states (id, updated_at, name, abbreviation) FROM stdin;
1	2017-03-14 22:14:32.213197	Alabama	AL
2	2017-03-14 22:14:32.213197	Alaska	AK
3	2017-03-14 22:14:32.213197	Arizona	AZ
4	2017-03-14 22:14:32.213197	Arkansas	AR
5	2017-03-14 22:14:32.213197	California	CA
6	2017-03-14 22:14:32.213197	Colorado	CO
7	2017-03-14 22:14:32.213197	Connecticut	CT
8	2017-03-14 22:14:32.213197	Delaware	DE
9	2017-03-14 22:14:32.213197	Florida	FL
10	2017-03-14 22:14:32.213197	Georgia	GA
11	2017-03-14 22:14:32.213197	Hawaii	HI
12	2017-03-14 22:14:32.213197	Idaho	ID
13	2017-03-14 22:14:32.213197	Illinois	IL
14	2017-03-14 22:14:32.213197	Indiana	IN
15	2017-03-14 22:14:32.213197	Iowa	IA
16	2017-03-14 22:14:32.213197	Kansas	KS
17	2017-03-14 22:14:32.213197	Kentucky	KY
18	2017-03-14 22:14:32.213197	Louisiana	LA
19	2017-03-14 22:14:32.213197	Maine	ME
20	2017-03-14 22:14:32.213197	Maryland	MD
21	2017-03-14 22:14:32.213197	Massachusetts	MA
22	2017-03-14 22:14:32.213197	Michigan	MI
23	2017-03-14 22:14:32.213197	Minnesota	MN
24	2017-03-14 22:14:32.213197	Mississippi	MS
25	2017-03-14 22:14:32.213197	Missouri	MO
26	2017-03-14 22:14:32.213197	Montana	MT
27	2017-03-14 22:14:32.213197	Nebraska	NE
28	2017-03-14 22:14:32.213197	Nevada	NV
29	2017-03-14 22:14:32.213197	New Hampshire	NH
30	2017-03-14 22:14:32.213197	New Jersey	NJ
31	2017-03-14 22:14:32.213197	New Mexico	NM
32	2017-03-14 22:14:32.213197	New York	NY
33	2017-03-14 22:14:32.213197	North Carolina	NC
34	2017-03-14 22:14:32.213197	North Dakota	ND
35	2017-03-14 22:14:32.213197	Ohio	OH
36	2017-03-14 22:14:32.213197	Oklahoma	OK
37	2017-03-14 22:14:32.213197	Oregon	OR
38	2017-03-14 22:14:32.213197	Pennsylvania	PA
39	2017-03-14 22:14:32.213197	Rhode Island	RI
40	2017-03-14 22:14:32.213197	South Carolina	SC
41	2017-03-14 22:14:32.213197	South Dakota	SD
42	2017-03-14 22:14:32.213197	Tennessee	TN
43	2017-03-14 22:14:32.213197	Texas	TX
44	2017-03-14 22:14:32.213197	Utah	UT
45	2017-03-14 22:14:32.213197	Vermont	VT
46	2017-03-14 22:14:32.213197	Virginia	VA
47	2017-03-14 22:14:32.213197	Washington	WA
48	2017-03-14 22:14:32.213197	West Virginia	WV
49	2017-03-14 22:14:32.213197	Wisconsin	WI
50	2017-03-14 22:14:32.213197	Wyoming	WY
51	2017-03-14 22:14:32.213197	Washington DC	DC
52	2017-03-14 22:14:32.213197	Puerto Rico	PR
53	2017-03-14 22:14:32.213197	U.S. Virgin Islands	VI
54	2017-03-14 22:14:32.213197	American Samoa	AS
55	2017-03-14 22:14:32.213197	Guam	GU
56	2017-03-14 22:14:32.213197	Northern Mariana Islands	MP
\.


--
-- TOC entry 2670 (class 0 OID 0)
-- Dependencies: 212
-- Name: states_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('states_id_seq', 1, false);


--
-- TOC entry 2612 (class 0 OID 33786)
-- Dependencies: 215
-- Data for Name: user_contacts; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY user_contacts (id, updated_at, user_id, people_id) FROM stdin;
1	2017-03-14 22:20:23.898557	1	7
2	2017-03-14 22:20:33.810175	1	1
3	2017-03-14 22:21:11.163938	1	2
4	2017-03-14 22:21:11.163938	1	3
5	2017-03-14 22:21:11.163938	1	5
6	2017-03-14 22:21:11.163938	1	4
7	2017-03-14 22:21:11.163938	1	6
8	2017-03-14 22:21:11.163938	1	8
9	2017-03-14 22:21:11.163938	1	9
10	2017-03-14 22:21:11.163938	1	10
\.


--
-- TOC entry 2671 (class 0 OID 0)
-- Dependencies: 214
-- Name: user_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('user_contacts_id_seq', 10, true);


--
-- TOC entry 2614 (class 0 OID 33795)
-- Dependencies: 217
-- Data for Name: user_contracts; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY user_contracts (id, updated_at, user_id, contracts_id) FROM stdin;
1	2017-03-14 22:48:28.38995	1	4
2	2017-03-14 22:48:28.38995	1	3
3	2017-03-14 22:48:28.38995	1	2
4	2017-03-14 22:48:28.38995	1	1
5	2017-03-14 22:48:28.38995	1	8
\.


--
-- TOC entry 2672 (class 0 OID 0)
-- Dependencies: 216
-- Name: user_contracts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('user_contracts_id_seq', 5, true);


--
-- TOC entry 2616 (class 0 OID 33804)
-- Dependencies: 219
-- Data for Name: user_facilities; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY user_facilities (id, updated_at, user_id, facility_id) FROM stdin;
1	2017-03-14 22:48:54.15319	1	4
2	2017-03-14 22:48:54.15319	1	3
3	2017-03-14 22:48:54.15319	1	2
4	2017-03-14 22:48:54.15319	1	1
5	2017-03-14 22:48:54.15319	1	6
\.


--
-- TOC entry 2673 (class 0 OID 0)
-- Dependencies: 218
-- Name: user_facilities_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('user_facilities_id_seq', 5, true);


--
-- TOC entry 2618 (class 0 OID 33813)
-- Dependencies: 221
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: chrislewis
--

COPY users (id, updated_at, people_id, email, password) FROM stdin;
1	2017-03-14 22:03:50.171613	\N	chrislewispac@gmail.com	$2a$10$VbpQkbEbczz3lsEqgBl7t.2r1Z.2XKIxu/1iLQQR7nNRFA6TI1z92
\.


--
-- TOC entry 2674 (class 0 OID 0)
-- Dependencies: 220
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: chrislewis
--

SELECT pg_catalog.setval('users_id_seq', 1, true);


--
-- TOC entry 2368 (class 2606 OID 33555)
-- Name: contact_instance contact_instance_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contact_instance
    ADD CONSTRAINT contact_instance_pk PRIMARY KEY (id);


--
-- TOC entry 2370 (class 2606 OID 33564)
-- Name: contract_facilities contract_facilities_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_facilities
    ADD CONSTRAINT contract_facilities_pk PRIMARY KEY (id);


--
-- TOC entry 2372 (class 2606 OID 33573)
-- Name: contract_status contract_status_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_status
    ADD CONSTRAINT contract_status_pk PRIMARY KEY (id);


--
-- TOC entry 2374 (class 2606 OID 33585)
-- Name: contracts contracts_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contracts
    ADD CONSTRAINT contracts_pk PRIMARY KEY (id);


--
-- TOC entry 2376 (class 2606 OID 33597)
-- Name: emrs emrs_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY emrs
    ADD CONSTRAINT emrs_pk PRIMARY KEY (id);


--
-- TOC entry 2378 (class 2606 OID 33610)
-- Name: facilities facilities_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_pk PRIMARY KEY (id);


--
-- TOC entry 2380 (class 2606 OID 33616)
-- Name: facility_area_hobbies facility_area_hobbies_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_area_hobbies
    ADD CONSTRAINT facility_area_hobbies_pk PRIMARY KEY (id);


--
-- TOC entry 2382 (class 2606 OID 33639)
-- Name: facility_stats facility_stats_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_stats
    ADD CONSTRAINT facility_stats_pk PRIMARY KEY (id);


--
-- TOC entry 2384 (class 2606 OID 33651)
-- Name: facility_types facility_types_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_types
    ADD CONSTRAINT facility_types_pk PRIMARY KEY (id);


--
-- TOC entry 2386 (class 2606 OID 34035)
-- Name: hobbies hobbies_name_key; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hobbies
    ADD CONSTRAINT hobbies_name_key UNIQUE (name);


--
-- TOC entry 2388 (class 2606 OID 33663)
-- Name: hobbies hobbies_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hobbies
    ADD CONSTRAINT hobbies_pk PRIMARY KEY (id);


--
-- TOC entry 2390 (class 2606 OID 33675)
-- Name: hospitalists hospitalists_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hospitalists
    ADD CONSTRAINT hospitalists_pk PRIMARY KEY (id);


--
-- TOC entry 2392 (class 2606 OID 33687)
-- Name: people people_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY people
    ADD CONSTRAINT people_pk PRIMARY KEY (id);


--
-- TOC entry 2394 (class 2606 OID 33699)
-- Name: provider_certifications provider_certifications_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_certifications
    ADD CONSTRAINT provider_certifications_pk PRIMARY KEY (id);


--
-- TOC entry 2396 (class 2606 OID 33711)
-- Name: provider_contracts provider_contracts_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_contracts
    ADD CONSTRAINT provider_contracts_pk PRIMARY KEY (id);


--
-- TOC entry 2398 (class 2606 OID 33720)
-- Name: provider_facilities provider_facilities_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_facilities
    ADD CONSTRAINT provider_facilities_pk PRIMARY KEY (id);


--
-- TOC entry 2400 (class 2606 OID 33729)
-- Name: provider_hobbies provider_hobbies_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_hobbies
    ADD CONSTRAINT provider_hobbies_pk PRIMARY KEY (id);


--
-- TOC entry 2402 (class 2606 OID 33741)
-- Name: provider_licenses provider_licenses_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_licenses
    ADD CONSTRAINT provider_licenses_pk PRIMARY KEY (id);


--
-- TOC entry 2404 (class 2606 OID 33750)
-- Name: provider_status provider_status_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_status
    ADD CONSTRAINT provider_status_pk PRIMARY KEY (id);


--
-- TOC entry 2406 (class 2606 OID 33762)
-- Name: providers providers_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY providers
    ADD CONSTRAINT providers_pk PRIMARY KEY (id);


--
-- TOC entry 2408 (class 2606 OID 33771)
-- Name: recruiters recruiters_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY recruiters
    ADD CONSTRAINT recruiters_pk PRIMARY KEY (id);


--
-- TOC entry 2410 (class 2606 OID 33783)
-- Name: states states_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY states
    ADD CONSTRAINT states_pk PRIMARY KEY (id);


--
-- TOC entry 2412 (class 2606 OID 33792)
-- Name: user_contacts user_contacts_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contacts
    ADD CONSTRAINT user_contacts_pk PRIMARY KEY (id);


--
-- TOC entry 2414 (class 2606 OID 33801)
-- Name: user_contracts user_contracts_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contracts
    ADD CONSTRAINT user_contracts_pk PRIMARY KEY (id);


--
-- TOC entry 2416 (class 2606 OID 33810)
-- Name: user_facilities user_facilities_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_facilities
    ADD CONSTRAINT user_facilities_pk PRIMARY KEY (id);


--
-- TOC entry 2418 (class 2606 OID 33822)
-- Name: users users_pk; Type: CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pk PRIMARY KEY (id);


--
-- TOC entry 2419 (class 2606 OID 34051)
-- Name: contact_instance contact_id; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contact_instance
    ADD CONSTRAINT contact_id FOREIGN KEY (contact_id) REFERENCES people(id);


--
-- TOC entry 2421 (class 2606 OID 34061)
-- Name: contact_instance contact_instance_user_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contact_instance
    ADD CONSTRAINT contact_instance_user_contract_id_fkey FOREIGN KEY (user_contract_id) REFERENCES contracts(id);


--
-- TOC entry 2420 (class 2606 OID 34056)
-- Name: contact_instance contact_instance_user_facility_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contact_instance
    ADD CONSTRAINT contact_instance_user_facility_id_fkey FOREIGN KEY (user_facility_id) REFERENCES facilities(id);


--
-- TOC entry 2423 (class 2606 OID 34041)
-- Name: contract_facilities contract_facilities_contracts_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_facilities
    ADD CONSTRAINT contract_facilities_contracts_id_fkey FOREIGN KEY (contracts_id) REFERENCES contracts(id);


--
-- TOC entry 2422 (class 2606 OID 34036)
-- Name: contract_facilities contract_facilities_facilities_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_facilities
    ADD CONSTRAINT contract_facilities_facilities_id_fkey FOREIGN KEY (facilities_id) REFERENCES facilities(id);


--
-- TOC entry 2424 (class 2606 OID 33848)
-- Name: contract_status contract_status_contracts; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contract_status
    ADD CONSTRAINT contract_status_contracts FOREIGN KEY (contract_id) REFERENCES contracts(id);


--
-- TOC entry 2425 (class 2606 OID 33853)
-- Name: contracts contracts_contact_instance; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contracts
    ADD CONSTRAINT contracts_contact_instance FOREIGN KEY (last_contact_id) REFERENCES contact_instance(id);


--
-- TOC entry 2426 (class 2606 OID 34046)
-- Name: contracts contracts_lead_source_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY contracts
    ADD CONSTRAINT contracts_lead_source_id_fkey FOREIGN KEY (lead_source_id) REFERENCES people(id);


--
-- TOC entry 2427 (class 2606 OID 33863)
-- Name: facilities facilities_contact_instance; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_contact_instance FOREIGN KEY (last_contact_id) REFERENCES contact_instance(id);


--
-- TOC entry 2428 (class 2606 OID 33868)
-- Name: facilities facilities_emrs; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_emrs FOREIGN KEY (emrs_id) REFERENCES emrs(id);


--
-- TOC entry 2429 (class 2606 OID 33873)
-- Name: facilities facilities_facility_stats; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_facility_stats FOREIGN KEY (facility_stats_id) REFERENCES facility_stats(id);


--
-- TOC entry 2430 (class 2606 OID 33878)
-- Name: facilities facilities_facility_types; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_facility_types FOREIGN KEY (facility_type_id) REFERENCES facility_types(id);


--
-- TOC entry 2431 (class 2606 OID 33883)
-- Name: facilities facilities_hospitalists; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_hospitalists FOREIGN KEY (hospitalist_id) REFERENCES hospitalists(id);


--
-- TOC entry 2432 (class 2606 OID 33888)
-- Name: facilities facilities_states; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facilities_states FOREIGN KEY (state_id) REFERENCES states(id);


--
-- TOC entry 2436 (class 2606 OID 33898)
-- Name: facility_area_hobbies facility_area_hobbies_facilities; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_area_hobbies
    ADD CONSTRAINT facility_area_hobbies_facilities FOREIGN KEY (facility_id) REFERENCES facilities(id);


--
-- TOC entry 2435 (class 2606 OID 33893)
-- Name: facility_area_hobbies facility_area_hobbies_hobbies; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facility_area_hobbies
    ADD CONSTRAINT facility_area_hobbies_hobbies FOREIGN KEY (hobby_id) REFERENCES hobbies(id);


--
-- TOC entry 2433 (class 2606 OID 33903)
-- Name: facilities facility_primary_contact; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facility_primary_contact FOREIGN KEY (primary_contact_id) REFERENCES user_contacts(id);


--
-- TOC entry 2434 (class 2606 OID 33908)
-- Name: facilities facility_secondary_contact; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY facilities
    ADD CONSTRAINT facility_secondary_contact FOREIGN KEY (secondary_contact_id) REFERENCES user_contacts(id);


--
-- TOC entry 2437 (class 2606 OID 33913)
-- Name: hospitalists hospitalists_user_contacts; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY hospitalists
    ADD CONSTRAINT hospitalists_user_contacts FOREIGN KEY (primary_contact_id) REFERENCES user_contacts(id);


--
-- TOC entry 2438 (class 2606 OID 33918)
-- Name: people people_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY people
    ADD CONSTRAINT people_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2439 (class 2606 OID 33923)
-- Name: people people_recruiters; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY people
    ADD CONSTRAINT people_recruiters FOREIGN KEY (recruiter_id) REFERENCES recruiters(id);


--
-- TOC entry 2440 (class 2606 OID 33928)
-- Name: people people_states; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY people
    ADD CONSTRAINT people_states FOREIGN KEY (state_id) REFERENCES states(id);


--
-- TOC entry 2441 (class 2606 OID 33933)
-- Name: provider_certifications provider_certifications_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_certifications
    ADD CONSTRAINT provider_certifications_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2442 (class 2606 OID 33938)
-- Name: provider_contracts provider_contracts_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_contracts
    ADD CONSTRAINT provider_contracts_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2443 (class 2606 OID 33943)
-- Name: provider_contracts provider_contracts_user_contracts; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_contracts
    ADD CONSTRAINT provider_contracts_user_contracts FOREIGN KEY (user_contracts_id) REFERENCES user_contracts(id);


--
-- TOC entry 2444 (class 2606 OID 33948)
-- Name: provider_contracts provider_contracts_user_facilities; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_contracts
    ADD CONSTRAINT provider_contracts_user_facilities FOREIGN KEY (user_facilities_id) REFERENCES user_facilities(id);


--
-- TOC entry 2445 (class 2606 OID 33953)
-- Name: provider_facilities provider_facilities_facilities; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_facilities
    ADD CONSTRAINT provider_facilities_facilities FOREIGN KEY (facility_id) REFERENCES facilities(id);


--
-- TOC entry 2446 (class 2606 OID 33958)
-- Name: provider_facilities provider_facilities_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_facilities
    ADD CONSTRAINT provider_facilities_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2447 (class 2606 OID 33963)
-- Name: provider_hobbies provider_hobbies_hobbies; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_hobbies
    ADD CONSTRAINT provider_hobbies_hobbies FOREIGN KEY (hobby_id) REFERENCES hobbies(id);


--
-- TOC entry 2448 (class 2606 OID 33968)
-- Name: provider_hobbies provider_hobbies_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_hobbies
    ADD CONSTRAINT provider_hobbies_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2449 (class 2606 OID 33973)
-- Name: provider_licenses provider_licenses_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_licenses
    ADD CONSTRAINT provider_licenses_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2450 (class 2606 OID 33978)
-- Name: provider_licenses provider_licenses_states; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_licenses
    ADD CONSTRAINT provider_licenses_states FOREIGN KEY (state_id) REFERENCES states(id);


--
-- TOC entry 2451 (class 2606 OID 33983)
-- Name: provider_status provider_status_providers; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY provider_status
    ADD CONSTRAINT provider_status_providers FOREIGN KEY (provider_id) REFERENCES providers(id);


--
-- TOC entry 2452 (class 2606 OID 33988)
-- Name: providers providers_people; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY providers
    ADD CONSTRAINT providers_people FOREIGN KEY (lead_source_id) REFERENCES people(id);


--
-- TOC entry 2453 (class 2606 OID 33993)
-- Name: providers providers_recruiters; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY providers
    ADD CONSTRAINT providers_recruiters FOREIGN KEY (responsible_recruiter_id) REFERENCES recruiters(id);


--
-- TOC entry 2454 (class 2606 OID 33998)
-- Name: user_contacts user_contacts_people; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contacts
    ADD CONSTRAINT user_contacts_people FOREIGN KEY (people_id) REFERENCES people(id);


--
-- TOC entry 2455 (class 2606 OID 34003)
-- Name: user_contacts user_contacts_users; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contacts
    ADD CONSTRAINT user_contacts_users FOREIGN KEY (user_id) REFERENCES users(id);


--
-- TOC entry 2456 (class 2606 OID 34008)
-- Name: user_contracts user_contracts_contracts; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contracts
    ADD CONSTRAINT user_contracts_contracts FOREIGN KEY (contracts_id) REFERENCES contracts(id);


--
-- TOC entry 2457 (class 2606 OID 34013)
-- Name: user_contracts user_contracts_users; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_contracts
    ADD CONSTRAINT user_contracts_users FOREIGN KEY (user_id) REFERENCES users(id);


--
-- TOC entry 2458 (class 2606 OID 34018)
-- Name: user_facilities user_facilities_facilities; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_facilities
    ADD CONSTRAINT user_facilities_facilities FOREIGN KEY (facility_id) REFERENCES facilities(id);


--
-- TOC entry 2459 (class 2606 OID 34023)
-- Name: user_facilities user_facilities_users; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY user_facilities
    ADD CONSTRAINT user_facilities_users FOREIGN KEY (user_id) REFERENCES users(id);


--
-- TOC entry 2460 (class 2606 OID 34028)
-- Name: users users_people; Type: FK CONSTRAINT; Schema: public; Owner: chrislewis
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_people FOREIGN KEY (people_id) REFERENCES people(id);


--
-- TOC entry 2625 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: chrislewis
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM chrislewis;
GRANT ALL ON SCHEMA public TO chrislewis;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2017-03-15 09:32:06 EDT

--
-- PostgreSQL database dump complete
--

