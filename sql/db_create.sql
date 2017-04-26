-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2017-03-16 04:09:27.692

create extension "uuid-ossp";

-- tables
-- Table: contact_instance
CREATE TABLE contact_instance (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    contact_type text  NULL,
    user_contacts_id text  NULL,
    contract_facilities_id text  NULL,
    user_facilities_id text  NULL,
    CONSTRAINT contact_instance_pk PRIMARY KEY (id)
);

-- Table: contract_facilities
CREATE TABLE contract_facilities (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    user_facilities_id text  NOT NULL,
    contracts_id text  NOT NULL,
    CONSTRAINT contract_facilities_pk PRIMARY KEY (id)
);

-- Table: contract_status
CREATE TABLE contract_status (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    lead_no_contact timestamp  NULL DEFAULT NULL,
    initial_contact timestamp  NULL,
    interested timestamp  NULL,
    not_interested_atm timestamp  NULL,
    proposal_sent timestamp  NULL,
    contract_sent timestamp  NULL,
    contract_signed timestamp  NULL,
    working timestamp  NULL,
    dead timestamp  NULL,
    undesirable timestamp  NULL,
    contracts_id text  NOT NULL,
    CONSTRAINT contract_status_pk PRIMARY KEY (id)
);

-- Table: contracts
CREATE TABLE contracts (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    renewal_date timestamp  NULL,
    file text  NULL,
    CONSTRAINT contracts_pk PRIMARY KEY (id)
);

-- Table: emrs
CREATE TABLE emrs (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NULL,
    CONSTRAINT emrs_pk PRIMARY KEY (id)
);

-- Table: facilities
CREATE TABLE facilities (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NULL,
    address text  NULL,
    city text  NULL,
    zip text  NULL,
    lat float  NULL,
    lon float  NULL,
    admission_orders bool  NULL DEFAULT FALSE,
    facility_type_id text  NULL,
    emr_id text  NULL,
    hospitalists_id text  NULL,
    states_id text  NULL,
    primary_contact_id text  NULL,
    secondary_contact_id text  NULL,
    CONSTRAINT facilities_pk PRIMARY KEY (id)
);

-- Table: facility_area_hobbies
CREATE TABLE facility_area_hobbies (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    hobbies_id text  NOT NULL,
    facilities_id text  NOT NULL,
    CONSTRAINT facility_area_hobbies_pk PRIMARY KEY (id)
);

-- Table: facility_stats
CREATE TABLE facility_stats (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    entry_date timestamp  NOT NULL DEFAULT now(),
    bed_count_ft int  NULL DEFAULT 0,
    bed_count_main int  NULL DEFAULT 0,
    bed_count_overflow int  NULL,
    annual_volume int  NULL DEFAULT 0,
    provider_in_triage bool  NULL DEFAULT FALSE,
    hours_physician int  NULL DEFAULT 0,
    hours_apc int  NULL DEFAULT 0,
    shifts_physician int  NULL DEFAULT 0,
    shifts_apc int  NULL DEFAULT 0,
    annual_hospitalizations int  NULL DEFAULT 0,
    annual_obs int  NULL DEFAULT 0,
    annual_admit int  NULL DEFAULT 0,
    lwbs int  NULL DEFAULT 0,
    lwot int  NULL DEFAULT 0,
    ft_hours_start timestamp  NULL,
    ft_hours_end timestamp  NULL,
    facility_id text  NOT NULL,
    CONSTRAINT facility_stats_pk PRIMARY KEY (id)
);

-- Table: facility_types
CREATE TABLE facility_types (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NOT NULL,
    CONSTRAINT facility_types_pk PRIMARY KEY (id)
);

-- Table: hobbies
CREATE TABLE hobbies (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NOT NULL,
    CONSTRAINT hobbies_pk PRIMARY KEY (id)
);

-- Table: hospitalists
CREATE TABLE hospitalists (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NULL,
    CONSTRAINT hospitalists_pk PRIMARY KEY (id)
);

-- Table: people
CREATE TABLE people (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    first_name text  NOT NULL,
    last_name text  NOT NULL,
    middle_name text  NULL,
    title text  NULL,
    email text  NOT NULL,
    cell_phone text  NULL,
    home_phone text  NULL,
    office_phone text  NULL,
    fax text  NULL,
    address text  NULL,
    zip text  NULL,
    city text  NULL,
    lat float  NULL,
    lon float  NULL,
    providers_id text  NOT NULL,
    recruiters_id text  NOT NULL,
    states_id text  NOT NULL,
    CONSTRAINT people_pk PRIMARY KEY (id)
);

-- Table: provider_certifications
CREATE TABLE provider_certifications (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NOT NULL,
    expiration_date timestamp  NOT NULL,
    file text  NULL,
    provider_id text  NOT NULL,
    CONSTRAINT provider_certifications_pk PRIMARY KEY (id)
);

-- Table: provider_contracts
CREATE TABLE provider_contracts (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    start_date timestamp  NULL,
    renewal_date timestamp  NULL,
    file text  NULL,
    facilities_id text  NULL,
    providers_id text  NOT NULL,
    CONSTRAINT provider_contracts_pk PRIMARY KEY (id)
);

-- Table: provider_facilities
CREATE TABLE provider_facilities (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    credentialed_date timestamp  NULL,
    renewal_date timestamp  NULL,
    providers_id text  NOT NULL,
    facilities_id text  NOT NULL,
    CONSTRAINT provider_facilities_pk PRIMARY KEY (id)
);

-- Table: provider_hobbies
CREATE TABLE provider_hobbies (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    hobbies_id text  NOT NULL,
    providers_id text  NOT NULL,
    CONSTRAINT provider_hobbies_pk PRIMARY KEY (id)
);

-- Table: provider_licenses
CREATE TABLE provider_licenses (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    license_number int  NOT NULL,
    file text  NULL,
    provider_id text  NOT NULL,
    CONSTRAINT provider_licenses_pk PRIMARY KEY (id)
);

-- Table: provider_status
CREATE TABLE provider_status (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    lead_no_contact timestamp  NULL,
    initial_contact timestamp  NULL,
    interested timestamp  NULL,
    not_interested_atm timestamp  NULL,
    verbal_commitment timestamp  NULL,
    contract_sent timestamp  NULL,
    contract_signed timestamp  NULL,
    dead timestamp  NULL,
    do_not_use timestamp  NULL,
    providers_id text  NOT NULL,
    users_id text  NOT NULL,
    CONSTRAINT provider_status_pk PRIMARY KEY (id)
);

-- Table: providers
CREATE TABLE providers (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    npi int  NULL,
    w9 text  NULL,
    direct_deposit_form text  NULL,
    hourly_rate int  NULL,
    desired_shifts_month int  NULL,
    max_shifts_month int  NULL,
    min_shifts_month int  NULL,
    full_time bool  NULL,
    part_time bool  NULL,
    prn bool  NULL,
    retired bool  NULL,
    notes text  NULL,
    insurance_certificate text  NULL,
    tb_expiration timestamp  NULL,
    tb_file text  NULL,
    flu_expiration timestamp  NULL,
    flu_file text  NULL,
    CONSTRAINT providers_pk PRIMARY KEY (id)
);

-- Table: recruiters
CREATE TABLE recruiters (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    CONSTRAINT recruiters_pk PRIMARY KEY (id)
);

-- Table: states
CREATE TABLE states (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    name text  NOT NULL,
    abbreviation text  NOT NULL,
    CONSTRAINT states_pk PRIMARY KEY (id)
);

-- Table: user_contacts
CREATE TABLE user_contacts (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    people_id text  NOT NULL,
    users_id text  NOT NULL,
    CONSTRAINT user_contacts_pk PRIMARY KEY (id)
);

-- Table: user_contracts
CREATE TABLE user_contracts (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    users_id text  NOT NULL,
    contracts_id text  NOT NULL,
    CONSTRAINT user_contracts_pk PRIMARY KEY (id)
);

-- Table: user_facilities
CREATE TABLE user_facilities (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    users_id text  NOT NULL,
    facilities_id text  NOT NULL,
    CONSTRAINT user_facilities_pk PRIMARY KEY (id)
);

-- Table: users
CREATE TABLE users (
    id text  NOT NULL DEFAULT uuid_generate_v4(),
    updated_at timestamp  NOT NULL DEFAULT now(),
    email text  NOT NULL,
    password text  NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

-- foreign keys
-- Reference: contact_instance_contract_facilities (table: contact_instance)
ALTER TABLE contact_instance ADD CONSTRAINT contact_instance_contract_facilities
    FOREIGN KEY (contract_facilities_id)
    REFERENCES contract_facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: contact_instance_user_contacts (table: contact_instance)
ALTER TABLE contact_instance ADD CONSTRAINT contact_instance_user_contacts
    FOREIGN KEY (user_contacts_id)
    REFERENCES user_contacts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: contact_instance_user_facilities (table: contact_instance)
ALTER TABLE contact_instance ADD CONSTRAINT contact_instance_user_facilities
    FOREIGN KEY (user_facilities_id)
    REFERENCES user_facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: contract_facilities_contracts (table: contract_facilities)
ALTER TABLE contract_facilities ADD CONSTRAINT contract_facilities_contracts
    FOREIGN KEY (contracts_id)
    REFERENCES contracts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: contract_facilities_user_facilities (table: contract_facilities)
ALTER TABLE contract_facilities ADD CONSTRAINT contract_facilities_user_facilities
    FOREIGN KEY (user_facilities_id)
    REFERENCES user_facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: contract_status_contracts (table: contract_status)
ALTER TABLE contract_status ADD CONSTRAINT contract_status_contracts
    FOREIGN KEY (contracts_id)
    REFERENCES contracts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facilities_emrs (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT facilities_emrs
    FOREIGN KEY (emr_id)
    REFERENCES emrs (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facilities_facility_types (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT facilities_facility_types
    FOREIGN KEY (facility_type_id)
    REFERENCES facility_types (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facilities_hospitalists (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT facilities_hospitalists
    FOREIGN KEY (hospitalists_id)
    REFERENCES hospitalists (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facilities_states (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT facilities_states
    FOREIGN KEY (states_id)
    REFERENCES states (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facilities_user_contacts (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT facilities_user_contacts
    FOREIGN KEY (secondary_contact_id)
    REFERENCES user_contacts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facility_area_hobbies_facilities (table: facility_area_hobbies)
ALTER TABLE facility_area_hobbies ADD CONSTRAINT facility_area_hobbies_facilities
    FOREIGN KEY (facilities_id)
    REFERENCES facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facility_area_hobbies_hobbies (table: facility_area_hobbies)
ALTER TABLE facility_area_hobbies ADD CONSTRAINT facility_area_hobbies_hobbies
    FOREIGN KEY (hobbies_id)
    REFERENCES hobbies (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: facility_stats_facilities (table: facility_stats)
ALTER TABLE facility_stats ADD CONSTRAINT facility_stats_facilities
    FOREIGN KEY (facility_id)
    REFERENCES facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: people_providers (table: people)
ALTER TABLE people ADD CONSTRAINT people_providers
    FOREIGN KEY (providers_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: people_recruiters (table: people)
ALTER TABLE people ADD CONSTRAINT people_recruiters
    FOREIGN KEY (recruiters_id)
    REFERENCES recruiters (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: people_states (table: people)
ALTER TABLE people ADD CONSTRAINT people_states
    FOREIGN KEY (states_id)
    REFERENCES states (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_certifications_providers (table: provider_certifications)
ALTER TABLE provider_certifications ADD CONSTRAINT provider_certifications_providers
    FOREIGN KEY (provider_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_contracts_facilities (table: provider_contracts)
ALTER TABLE provider_contracts ADD CONSTRAINT provider_contracts_facilities
    FOREIGN KEY (facilities_id)
    REFERENCES facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_contracts_providers (table: provider_contracts)
ALTER TABLE provider_contracts ADD CONSTRAINT provider_contracts_providers
    FOREIGN KEY (providers_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_facilities_facilities (table: provider_facilities)
ALTER TABLE provider_facilities ADD CONSTRAINT provider_facilities_facilities
    FOREIGN KEY (facilities_id)
    REFERENCES facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_facilities_providers (table: provider_facilities)
ALTER TABLE provider_facilities ADD CONSTRAINT provider_facilities_providers
    FOREIGN KEY (providers_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_hobbies_hobbies (table: provider_hobbies)
ALTER TABLE provider_hobbies ADD CONSTRAINT provider_hobbies_hobbies
    FOREIGN KEY (hobbies_id)
    REFERENCES hobbies (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_hobbies_providers (table: provider_hobbies)
ALTER TABLE provider_hobbies ADD CONSTRAINT provider_hobbies_providers
    FOREIGN KEY (providers_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_licenses_providers (table: provider_licenses)
ALTER TABLE provider_licenses ADD CONSTRAINT provider_licenses_providers
    FOREIGN KEY (provider_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_status_providers (table: provider_status)
ALTER TABLE provider_status ADD CONSTRAINT provider_status_providers
    FOREIGN KEY (providers_id)
    REFERENCES providers (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: provider_status_users (table: provider_status)
ALTER TABLE provider_status ADD CONSTRAINT provider_status_users
    FOREIGN KEY (users_id)
    REFERENCES users (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: secondary_contact (table: facilities)
ALTER TABLE facilities ADD CONSTRAINT secondary_contact
    FOREIGN KEY (primary_contact_id)
    REFERENCES user_contacts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_contacts_people (table: user_contacts)
ALTER TABLE user_contacts ADD CONSTRAINT user_contacts_people
    FOREIGN KEY (people_id)
    REFERENCES people (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_contacts_users (table: user_contacts)
ALTER TABLE user_contacts ADD CONSTRAINT user_contacts_users
    FOREIGN KEY (users_id)
    REFERENCES users (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_contracts_contracts (table: user_contracts)
ALTER TABLE user_contracts ADD CONSTRAINT user_contracts_contracts
    FOREIGN KEY (contracts_id)
    REFERENCES contracts (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_contracts_users (table: user_contracts)
ALTER TABLE user_contracts ADD CONSTRAINT user_contracts_users
    FOREIGN KEY (users_id)
    REFERENCES users (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_facilities_facilities (table: user_facilities)
ALTER TABLE user_facilities ADD CONSTRAINT user_facilities_facilities
    FOREIGN KEY (facilities_id)
    REFERENCES facilities (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference: user_facilities_users (table: user_facilities)
ALTER TABLE user_facilities ADD CONSTRAINT user_facilities_users
    FOREIGN KEY (users_id)
    REFERENCES users (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- End of file.
