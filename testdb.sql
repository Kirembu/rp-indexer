DROP TABLE IF EXISTS contacts_contact CASCADE;
CREATE TABLE contacts_contact (
    id integer NOT NULL,
    is_active boolean NOT NULL,
    created_by_id integer NOT NULL,
    created_on timestamp with time zone NOT NULL,
    modified_by_id integer NOT NULL,
    modified_on timestamp with time zone NOT NULL,
    org_id integer NOT NULL,
    is_blocked boolean NOT NULL,
    name character varying(128),
    is_test boolean NOT NULL,
    language character varying(3),
    uuid character varying(36) NOT NULL,
    is_stopped boolean NOT NULL,
    fields jsonb
);

CREATE SEQUENCE contacts_contact_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contacts_contact_id_seq OWNED BY contacts_contact.id;

DROP TABLE IF EXISTS contacts_contacturn CASCADE;
CREATE TABLE contacts_contacturn (
    id integer NOT NULL,
    contact_id integer,
    scheme character varying(128) NOT NULL,
    org_id integer NOT NULL,
    priority integer NOT NULL,
    path character varying(255) NOT NULL,
    channel_id integer,
    auth text,
    display character varying(255),
    identity character varying(255) NOT NULL
);

CREATE SEQUENCE contacts_contacturn_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contacts_contacturn_id_seq OWNED BY contacts_contacturn.id;

DROP TABLE IF EXISTS contacts_contactgroup CASCADE;
CREATE TABLE contacts_contactgroup (
    id integer NOT NULL,
    uuid character varying(36) NOT NULL,
    name character varying(128) NOT NULL
);

CREATE SEQUENCE contacts_contactgroup_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contacts_contactgroup_id_seq OWNED BY contacts_contactgroup.id;

DROP TABLE IF EXISTS contacts_contactgroup_contacts CASCADE;
CREATE TABLE contacts_contactgroup_contacts (
    id integer NOT NULL,
    contactgroup_id integer NOT NULL,
    contact_id integer NOT NULL
);

CREATE SEQUENCE contacts_contactgroup_contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE contacts_contactgroup_contacts_id_seq OWNED BY contacts_contactgroup_contacts.id;

-- Fields:
-- 17103bb1-1b48-4b70-92f7-1f6b73bd3488 - nickname (text)
-- 05bca1cd-e322-4837-9595-86d0d85e5adb - age (numeric)
-- e0eac267-463a-4c00-9732-cab62df07b16 - joined_on (datetime)
-- 22d11697-edba-4186-b084-793e3b876379 - home_state (state)
-- fcab2439-861c-4832-aa54-0c97f38f24ab - home_district (district)
-- a551ade4-e5a0-4d83-b185-53b515ad2f2a - home_ward (ward)

INSERT INTO contacts_contact(id, is_active, created_by_id, created_on, modified_by_id, modified_on, org_id, is_blocked, name, is_test, language, uuid, is_stopped, fields) VALUES
(1,  TRUE, -1, '2017-11-10 21:11:59.890662+00', -1, '2017-11-10 21:11:59.890662+00', 1, FALSE, NULL, FALSE, 'eng', 'c7a2dd87-a80e-420b-8431-ca48d422e924', FALSE, 
'{ "17103bb1-1b48-4b70-92f7-1f6b73bd3488": {"text": "the rock"}}'),
(2,  TRUE, -1, '2015-03-25 17:15:12.982168+00', -1, '2015-03-25 17:15:12.982168+00', 1, FALSE, NULL, TRUE, 'fra', '1ad43adc-c4fc-4244-8b3d-a938b8eba57a', FALSE, NULL),
(3,  TRUE, -1, '2015-03-26 10:07:14.054521+00', -1, '2015-03-26 10:07:14.054521+00', 1, FALSE, NULL, FALSE, NULL, '7a6606c7-ff41-4203-aa98-454a10d37209', TRUE,
'{ "05bca1cd-e322-4837-9595-86d0d85e5adb": {"text": "11", "number": 11 }}'),
(4,  TRUE, -1, '2015-03-26 13:04:58.699648+00', -1, '2015-03-26 13:04:58.699648+00', 1, TRUE, NULL, FALSE, NULL, '29b45297-15ad-4061-a7d4-e0b33d121541', FALSE, 
'{ "05bca1cd-e322-4837-9595-86d0d85e5adb": {"text": "9", "number": 9 }, "e0eac267-463a-4c00-9732-cab62df07b16": { "text": "2018-04-06T18:37:59+00:00", "datetime": "2018-04-06T18:37:59+00:00"}}'),
(5,  TRUE, -1, '2015-03-27 07:39:28.955051+00', -1, '2015-03-27 07:39:28.955051+00', 1, FALSE, 'John Doe', FALSE, NULL, '51762bba-01a2-4c4e-b5cd-b182d0405cd4', FALSE, 
'{ "e0eac267-463a-4c00-9732-cab62df07b16": { "text": "2030-04-06T18:37:59+00:00", "datetime": "2030-04-06T18:37:59+00:00"}}'),
(6,  TRUE, -1, '2015-10-30 19:42:27.001837+00', -1, '2015-10-30 19:42:27.001837+00', 2, FALSE, 'Ajodinabiff Dane', FALSE, NULL, '3e814add-e614-41f7-8b5d-a07f670a698f', FALSE, 
'{ "22d11697-edba-4186-b084-793e3b876379": { "text": "USA > Washington", "state": "USA > Washington"} }'),
(7,  TRUE, -1, '2017-11-10 21:11:59.890662+00', -1, '2017-11-10 21:11:59.890662+00', 2, FALSE, 'Joanne Stone', FALSE, NULL, '7051dff0-0a27-49d7-af1f-4494239139e6', FALSE, 
'{ "22d11697-edba-4186-b084-793e3b876379": { "text": "USA > Colorado", "state": "USA > Colorado"} }'),
(8,  TRUE, -1, '2015-03-27 13:39:43.995812+00', -1, '2015-03-27 13:39:43.995812+00', 2, FALSE, NULL, FALSE, NULL, 'b46f6e18-95b4-4984-9926-dded047f4eb3', FALSE, 
'{ "fcab2439-861c-4832-aa54-0c97f38f24ab": { "text": "USA > Washington > King County", "district": "USA > Washington > King County"} }'),
(9,  TRUE, -1, '2017-11-10 21:11:59.890662+00', -1, '2017-11-10 21:11:59.890662+00', 2, FALSE, NULL, FALSE, NULL, '9195c8b7-6138-4d84-ac56-5192cc3d8ceb', FALSE, 
'{ "a551ade4-e5a0-4d83-b185-53b515ad2f2a": { "text": "USA > Washington > King County > Central District", "ward": "USA > Washington > King County > Central District"} }'),
(10, TRUE, -1, '2016-08-22 14:20:05.690311+00', -1, '2016-08-22 14:20:05.690311+00', 2, FALSE, NULL, FALSE, NULL, '2b8bd28d-43e0-4c34-a4bb-0f10b11fdb8a', FALSE, 
'{ "fcab2439-861c-4832-aa54-0c97f38f24ab": { "text": "USA > Colorado > King", "district": "USA > Colorado > King"} }');

INSERT INTO contacts_contacturn(id, contact_id, scheme, org_id, priority, path, display, identity) VALUES
(1, 1, 'tel', 1, 50, '+12067791111', NULL, 'tel:+12067791111'),
(2, 1, 'tel', 1, 50, '+12067792222', NULL, 'tel:+12067792222'),
(3, 2, 'tel', 1, 50, '+12067793333', NULL, 'tel:+12067793333'),
(4, 3, 'tel', 1, 50, '+12067794444', NULL, 'tel:+12067794444'),
(5, 4, 'tel', 1, 50, '+12067795555', NULL, 'tel:+12067795555'),
(6, 5, 'tel', 1, 50, '+12060000556', NULL, 'tel:+12067796666'),
(7, 6, 'tel', 2, 50, '+12060005577', NULL, 'tel:+12067797777'),
(8, 7, 'tel', 2, 50, '+12067798888', NULL, 'tel:+12067798888'),
(9, 8, 'viber', 2, 90, 'viberpath==', NULL, 'viber:viberpath=='),
(10, 9, 'facebook', 2, 90, 1000001, 'funguy', 'facebook:1000001'),
(11, 10, 'twitterid', 2, 90, 1000001, 'fungal', 'twitterid:1000001'),
(12, 11, 'whatsapp',  2, 90, 1000003, NULL, 'whatsapp:1000003');

INSERT INTO contacts_contactgroup(id, uuid, name) VALUES
(1, '4ea0f313-2f62-4e57-bdf0-232b5191dd57', 'Group 1'),
(2, '4c016340-468d-4675-a974-15cb7a45a5ab', 'Group 2'),
(3, 'e61b5bf7-8ddf-4e05-b0a8-4c46a6b68cff', 'Group 3'),
(4, '529bac39-550a-4d6f-817c-1833f3449007', 'Group 4');

INSERT INTO contacts_contactgroup_contacts(id, contact_id, contactgroup_id) VALUES
(1, 1, 1),
(2, 2, 1),
(3, 1, 4),
(4, 3, 4);
