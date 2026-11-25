--
-- PostgreSQL database dump
--

\restrict Qfr1msnuqV1y5noreQhpKMqrbZz3TL3UmSH3ShIIla2bLJho7SE9Jous6xMWUih

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg13+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: blog; Type: SCHEMA; Schema: -; Owner: web
--

CREATE SCHEMA blog;


ALTER SCHEMA blog OWNER TO web;

--
-- Name: SCHEMA blog; Type: COMMENT; Schema: -; Owner: web
--

COMMENT ON SCHEMA blog IS 'Schema "blog" for user "web" in database "web".';


--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'Schema "public" for user "web" in database "web".';


--
-- Name: ums; Type: SCHEMA; Schema: -; Owner: web
--

CREATE SCHEMA ums;


ALTER SCHEMA ums OWNER TO web;

--
-- Name: SCHEMA ums; Type: COMMENT; Schema: -; Owner: web
--

COMMENT ON SCHEMA ums IS 'Schema "ums" for user "web" in database "web".';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: on_update_current_timestamp_article(); Type: FUNCTION; Schema: blog; Owner: web
--

CREATE FUNCTION blog.on_update_current_timestamp_article() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;


ALTER FUNCTION blog.on_update_current_timestamp_article() OWNER TO web;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: access; Type: TABLE; Schema: blog; Owner: web
--

CREATE TABLE blog.access (
    id bigint NOT NULL,
    ipv4 bigint DEFAULT '0'::bigint,
    "timestamp" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ipv6 bytea,
    aid bigint NOT NULL
);


ALTER TABLE blog.access OWNER TO web;

--
-- Name: access_id_seq; Type: SEQUENCE; Schema: blog; Owner: web
--

CREATE SEQUENCE blog.access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE blog.access_id_seq OWNER TO web;

--
-- Name: access_id_seq; Type: SEQUENCE OWNED BY; Schema: blog; Owner: web
--

ALTER SEQUENCE blog.access_id_seq OWNED BY blog.access.id;


--
-- Name: article; Type: TABLE; Schema: blog; Owner: web
--

CREATE TABLE blog.article (
    aid bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    title character varying(65) NOT NULL,
    markdown text,
    summary text
);


ALTER TABLE blog.article OWNER TO web;

--
-- Name: article_aid_seq; Type: SEQUENCE; Schema: blog; Owner: web
--

CREATE SEQUENCE blog.article_aid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE blog.article_aid_seq OWNER TO web;

--
-- Name: article_aid_seq; Type: SEQUENCE OWNED BY; Schema: blog; Owner: web
--

ALTER SEQUENCE blog.article_aid_seq OWNED BY blog.article.aid;


--
-- Name: comment_id_seq; Type: SEQUENCE; Schema: blog; Owner: web
--

CREATE SEQUENCE blog.comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE blog.comment_id_seq OWNER TO web;

--
-- Name: comment; Type: TABLE; Schema: blog; Owner: web
--

CREATE TABLE blog.comment (
    id bigint DEFAULT nextval('blog.comment_id_seq'::regclass) NOT NULL,
    uid bigint NOT NULL,
    aid bigint NOT NULL,
    content text NOT NULL,
    parent bigint NOT NULL,
    ip inet,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_deleted boolean DEFAULT false NOT NULL,
    reply_to character varying(32)
);


ALTER TABLE blog.comment OWNER TO web;

--
-- Name: tag; Type: TABLE; Schema: blog; Owner: web
--

CREATE TABLE blog.tag (
    id bigint NOT NULL,
    name character varying(30) NOT NULL
);


ALTER TABLE blog.tag OWNER TO web;

--
-- Name: tag_article; Type: TABLE; Schema: blog; Owner: web
--

CREATE TABLE blog.tag_article (
    id bigint NOT NULL,
    tag_id bigint NOT NULL,
    aid bigint NOT NULL
);


ALTER TABLE blog.tag_article OWNER TO web;

--
-- Name: tag_article_id_seq; Type: SEQUENCE; Schema: blog; Owner: web
--

CREATE SEQUENCE blog.tag_article_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE blog.tag_article_id_seq OWNER TO web;

--
-- Name: tag_article_id_seq; Type: SEQUENCE OWNED BY; Schema: blog; Owner: web
--

ALTER SEQUENCE blog.tag_article_id_seq OWNED BY blog.tag_article.id;


--
-- Name: tag_id_seq; Type: SEQUENCE; Schema: blog; Owner: web
--

CREATE SEQUENCE blog.tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE blog.tag_id_seq OWNER TO web;

--
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: blog; Owner: web
--

ALTER SEQUENCE blog.tag_id_seq OWNED BY blog.tag.id;


--
-- Name: rbac_permission; Type: TABLE; Schema: ums; Owner: web
--

CREATE TABLE ums.rbac_permission (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    description text
);


ALTER TABLE ums.rbac_permission OWNER TO web;

--
-- Name: rbac_permission_id_seq; Type: SEQUENCE; Schema: ums; Owner: web
--

CREATE SEQUENCE ums.rbac_permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.rbac_permission_id_seq OWNER TO web;

--
-- Name: rbac_permission_id_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: web
--

ALTER SEQUENCE ums.rbac_permission_id_seq OWNED BY ums.rbac_permission.id;


--
-- Name: rbac_role; Type: TABLE; Schema: ums; Owner: web
--

CREATE TABLE ums.rbac_role (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    description text
);


ALTER TABLE ums.rbac_role OWNER TO web;

--
-- Name: rbac_role_id_seq; Type: SEQUENCE; Schema: ums; Owner: web
--

CREATE SEQUENCE ums.rbac_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.rbac_role_id_seq OWNER TO web;

--
-- Name: rbac_role_id_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: web
--

ALTER SEQUENCE ums.rbac_role_id_seq OWNED BY ums.rbac_role.id;


--
-- Name: rbac_role_permission; Type: TABLE; Schema: ums; Owner: web
--

CREATE TABLE ums.rbac_role_permission (
    id bigint NOT NULL,
    role_id integer NOT NULL,
    permission_id integer NOT NULL
);


ALTER TABLE ums.rbac_role_permission OWNER TO web;

--
-- Name: rbac_role_permission_id_seq; Type: SEQUENCE; Schema: ums; Owner: web
--

CREATE SEQUENCE ums.rbac_role_permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.rbac_role_permission_id_seq OWNER TO web;

--
-- Name: rbac_role_permission_id_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: web
--

ALTER SEQUENCE ums.rbac_role_permission_id_seq OWNED BY ums.rbac_role_permission.id;


--
-- Name: rbac_user; Type: TABLE; Schema: ums; Owner: web
--

CREATE TABLE ums.rbac_user (
    uid bigint NOT NULL,
    phone character varying(20),
    username character varying(32) NOT NULL,
    nickname character varying(100) DEFAULT 'User'::character varying NOT NULL,
    avatar character varying(255) DEFAULT '/static/public/images/user/avatar.jpg'::character varying NOT NULL,
    email character varying(30) NOT NULL,
    gender integer DEFAULT '-1'::integer NOT NULL,
    options character varying(3000) DEFAULT '{}'::character varying NOT NULL,
    password character varying(255) NOT NULL,
    lock integer DEFAULT 0 NOT NULL,
    register_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_login_at timestamp with time zone DEFAULT '2017-01-01 00:00:00+00'::timestamp with time zone NOT NULL
);


ALTER TABLE ums.rbac_user OWNER TO web;

--
-- Name: rbac_user_role; Type: TABLE; Schema: ums; Owner: web
--

CREATE TABLE ums.rbac_user_role (
    id bigint NOT NULL,
    uid bigint NOT NULL,
    role_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE ums.rbac_user_role OWNER TO web;

--
-- Name: rbac_user_role_id_seq; Type: SEQUENCE; Schema: ums; Owner: web
--

CREATE SEQUENCE ums.rbac_user_role_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.rbac_user_role_id_seq OWNER TO web;

--
-- Name: rbac_user_role_id_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: web
--

ALTER SEQUENCE ums.rbac_user_role_id_seq OWNED BY ums.rbac_user_role.id;


--
-- Name: user_login_audit; Type: TABLE; Schema: ums; Owner: postgres
--

CREATE TABLE ums.user_login_audit (
    id bigint NOT NULL,
    uid bigint NOT NULL,
    login_ip character varying(45) NOT NULL,
    device_info text NOT NULL,
    login_method character varying(32) NOT NULL,
    refresh_jti character varying(64) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE ums.user_login_audit OWNER TO postgres;

--
-- Name: user_login_audit_id_seq; Type: SEQUENCE; Schema: ums; Owner: postgres
--

CREATE SEQUENCE ums.user_login_audit_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.user_login_audit_id_seq OWNER TO postgres;

--
-- Name: user_login_audit_id_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: postgres
--

ALTER SEQUENCE ums.user_login_audit_id_seq OWNED BY ums.user_login_audit.id;


--
-- Name: user_uid_seq; Type: SEQUENCE; Schema: ums; Owner: web
--

CREATE SEQUENCE ums.user_uid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE ums.user_uid_seq OWNER TO web;

--
-- Name: user_uid_seq; Type: SEQUENCE OWNED BY; Schema: ums; Owner: web
--

ALTER SEQUENCE ums.user_uid_seq OWNED BY ums.rbac_user.uid;


--
-- Name: access id; Type: DEFAULT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.access ALTER COLUMN id SET DEFAULT nextval('blog.access_id_seq'::regclass);


--
-- Name: article aid; Type: DEFAULT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.article ALTER COLUMN aid SET DEFAULT nextval('blog.article_aid_seq'::regclass);


--
-- Name: tag id; Type: DEFAULT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag ALTER COLUMN id SET DEFAULT nextval('blog.tag_id_seq'::regclass);


--
-- Name: tag_article id; Type: DEFAULT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag_article ALTER COLUMN id SET DEFAULT nextval('blog.tag_article_id_seq'::regclass);


--
-- Name: rbac_permission id; Type: DEFAULT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_permission ALTER COLUMN id SET DEFAULT nextval('ums.rbac_permission_id_seq'::regclass);


--
-- Name: rbac_role id; Type: DEFAULT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role ALTER COLUMN id SET DEFAULT nextval('ums.rbac_role_id_seq'::regclass);


--
-- Name: rbac_role_permission id; Type: DEFAULT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role_permission ALTER COLUMN id SET DEFAULT nextval('ums.rbac_role_permission_id_seq'::regclass);


--
-- Name: rbac_user uid; Type: DEFAULT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user ALTER COLUMN uid SET DEFAULT nextval('ums.user_uid_seq'::regclass);


--
-- Name: rbac_user_role id; Type: DEFAULT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user_role ALTER COLUMN id SET DEFAULT nextval('ums.rbac_user_role_id_seq'::regclass);


--
-- Name: user_login_audit id; Type: DEFAULT; Schema: ums; Owner: postgres
--

ALTER TABLE ONLY ums.user_login_audit ALTER COLUMN id SET DEFAULT nextval('ums.user_login_audit_id_seq'::regclass);


--
-- Name: comment comment_pkey; Type: CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.comment
    ADD CONSTRAINT comment_pkey PRIMARY KEY (id);


--
-- Name: access idx_24668_primary; Type: CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.access
    ADD CONSTRAINT idx_24668_primary PRIMARY KEY (id);


--
-- Name: article idx_24677_primary; Type: CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.article
    ADD CONSTRAINT idx_24677_primary PRIMARY KEY (aid);


--
-- Name: tag idx_24686_primary; Type: CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag
    ADD CONSTRAINT idx_24686_primary PRIMARY KEY (id);


--
-- Name: tag_article idx_24691_primary; Type: CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag_article
    ADD CONSTRAINT idx_24691_primary PRIMARY KEY (id);


--
-- Name: rbac_permission idx_24759_primary; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_permission
    ADD CONSTRAINT idx_24759_primary PRIMARY KEY (id);


--
-- Name: rbac_role idx_24763_primary; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role
    ADD CONSTRAINT idx_24763_primary PRIMARY KEY (id);


--
-- Name: rbac_role_permission idx_24768_primary; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role_permission
    ADD CONSTRAINT idx_24768_primary PRIMARY KEY (id);


--
-- Name: rbac_user_role idx_24773_primary; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user_role
    ADD CONSTRAINT idx_24773_primary PRIMARY KEY (id);


--
-- Name: rbac_user idx_24779_primary; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user
    ADD CONSTRAINT idx_24779_primary PRIMARY KEY (uid);


--
-- Name: rbac_role_permission unique_role_permission; Type: CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role_permission
    ADD CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id);


--
-- Name: user_login_audit user_login_audit_pkey; Type: CONSTRAINT; Schema: ums; Owner: postgres
--

ALTER TABLE ONLY ums.user_login_audit
    ADD CONSTRAINT user_login_audit_pkey PRIMARY KEY (id);


--
-- Name: idx_24686_index_name; Type: INDEX; Schema: blog; Owner: web
--

CREATE UNIQUE INDEX idx_24686_index_name ON blog.tag USING btree (name);


--
-- Name: idx_24691_index_aid; Type: INDEX; Schema: blog; Owner: web
--

CREATE INDEX idx_24691_index_aid ON blog.tag_article USING btree (aid);


--
-- Name: idx_24691_index_tag_id; Type: INDEX; Schema: blog; Owner: web
--

CREATE INDEX idx_24691_index_tag_id ON blog.tag_article USING btree (tag_id);


--
-- Name: idx_24768_rbac_role_permission_ibfk_1; Type: INDEX; Schema: ums; Owner: web
--

CREATE INDEX idx_24768_rbac_role_permission_ibfk_1 ON ums.rbac_role_permission USING btree (role_id);


--
-- Name: idx_24768_rbac_role_permission_ibfk_2; Type: INDEX; Schema: ums; Owner: web
--

CREATE INDEX idx_24768_rbac_role_permission_ibfk_2 ON ums.rbac_role_permission USING btree (permission_id);


--
-- Name: idx_24773_rbac_user_role_ibfk_1; Type: INDEX; Schema: ums; Owner: web
--

CREATE INDEX idx_24773_rbac_user_role_ibfk_1 ON ums.rbac_user_role USING btree (uid);


--
-- Name: idx_24773_rbac_user_role_ibfk_2; Type: INDEX; Schema: ums; Owner: web
--

CREATE INDEX idx_24773_rbac_user_role_ibfk_2 ON ums.rbac_user_role USING btree (role_id);


--
-- Name: idx_24779_index_email; Type: INDEX; Schema: ums; Owner: web
--

CREATE UNIQUE INDEX idx_24779_index_email ON ums.rbac_user USING btree (email);


--
-- Name: idx_24779_index_phone; Type: INDEX; Schema: ums; Owner: web
--

CREATE UNIQUE INDEX idx_24779_index_phone ON ums.rbac_user USING btree (phone);


--
-- Name: idx_24779_index_username; Type: INDEX; Schema: ums; Owner: web
--

CREATE UNIQUE INDEX idx_24779_index_username ON ums.rbac_user USING btree (username);


--
-- Name: idx_user_login_audit_created_at; Type: INDEX; Schema: ums; Owner: postgres
--

CREATE INDEX idx_user_login_audit_created_at ON ums.user_login_audit USING btree (created_at DESC);


--
-- Name: idx_user_login_audit_uid_created_at; Type: INDEX; Schema: ums; Owner: postgres
--

CREATE INDEX idx_user_login_audit_uid_created_at ON ums.user_login_audit USING btree (uid, created_at DESC);


--
-- Name: article on_update_current_timestamp; Type: TRIGGER; Schema: blog; Owner: web
--

CREATE TRIGGER on_update_current_timestamp BEFORE UPDATE ON blog.article FOR EACH ROW EXECUTE FUNCTION blog.on_update_current_timestamp_article();


--
-- Name: comment comment_ibfk_1; Type: FK CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.comment
    ADD CONSTRAINT comment_ibfk_1 FOREIGN KEY (uid) REFERENCES ums.rbac_user(uid) ON UPDATE CASCADE;


--
-- Name: comment comment_ibfk_2; Type: FK CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.comment
    ADD CONSTRAINT comment_ibfk_2 FOREIGN KEY (aid) REFERENCES blog.article(aid) ON UPDATE CASCADE;


--
-- Name: tag_article tag_article_ibfk_1; Type: FK CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag_article
    ADD CONSTRAINT tag_article_ibfk_1 FOREIGN KEY (tag_id) REFERENCES blog.tag(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: tag_article tag_article_ibfk_2; Type: FK CONSTRAINT; Schema: blog; Owner: web
--

ALTER TABLE ONLY blog.tag_article
    ADD CONSTRAINT tag_article_ibfk_2 FOREIGN KEY (aid) REFERENCES blog.article(aid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rbac_role_permission rbac_role_permission_ibfk_1; Type: FK CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role_permission
    ADD CONSTRAINT rbac_role_permission_ibfk_1 FOREIGN KEY (role_id) REFERENCES ums.rbac_role(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rbac_role_permission rbac_role_permission_ibfk_2; Type: FK CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_role_permission
    ADD CONSTRAINT rbac_role_permission_ibfk_2 FOREIGN KEY (permission_id) REFERENCES ums.rbac_permission(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rbac_user_role rbac_user_role_ibfk_1; Type: FK CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user_role
    ADD CONSTRAINT rbac_user_role_ibfk_1 FOREIGN KEY (uid) REFERENCES ums.rbac_user(uid) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: rbac_user_role rbac_user_role_ibfk_2; Type: FK CONSTRAINT; Schema: ums; Owner: web
--

ALTER TABLE ONLY ums.rbac_user_role
    ADD CONSTRAINT rbac_user_role_ibfk_2 FOREIGN KEY (role_id) REFERENCES ums.rbac_role(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_login_audit user_login_audit_uid_fkey; Type: FK CONSTRAINT; Schema: ums; Owner: postgres
--

ALTER TABLE ONLY ums.user_login_audit
    ADD CONSTRAINT user_login_audit_uid_fkey FOREIGN KEY (uid) REFERENCES ums.rbac_user(uid) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO web;


--
-- Name: FUNCTION uuid_generate_v1(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v1() TO web;


--
-- Name: FUNCTION uuid_generate_v1mc(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v1mc() TO web;


--
-- Name: FUNCTION uuid_generate_v3(namespace uuid, name text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v3(namespace uuid, name text) TO web;


--
-- Name: FUNCTION uuid_generate_v4(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v4() TO web;


--
-- Name: FUNCTION uuid_generate_v5(namespace uuid, name text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v5(namespace uuid, name text) TO web;


--
-- Name: FUNCTION uuid_nil(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_nil() TO web;


--
-- Name: FUNCTION uuid_ns_dns(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_dns() TO web;


--
-- Name: FUNCTION uuid_ns_oid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_oid() TO web;


--
-- Name: FUNCTION uuid_ns_url(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_url() TO web;


--
-- Name: FUNCTION uuid_ns_x500(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_x500() TO web;


--
-- Name: TABLE user_login_audit; Type: ACL; Schema: ums; Owner: postgres
--

GRANT ALL ON TABLE ums.user_login_audit TO web;


--
-- Name: SEQUENCE user_login_audit_id_seq; Type: ACL; Schema: ums; Owner: postgres
--

GRANT ALL ON SEQUENCE ums.user_login_audit_id_seq TO web;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: blog; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA blog GRANT ALL ON SEQUENCES TO web;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: blog; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA blog GRANT ALL ON FUNCTIONS TO web;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: blog; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA blog GRANT ALL ON TABLES TO web;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON SEQUENCES TO web;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON FUNCTIONS TO web;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT ALL ON TABLES TO web;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: ums; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA ums GRANT ALL ON SEQUENCES TO web;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: ums; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA ums GRANT ALL ON FUNCTIONS TO web;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: ums; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA ums GRANT ALL ON TABLES TO web;


--
-- PostgreSQL database dump complete
--

\unrestrict Qfr1msnuqV1y5noreQhpKMqrbZz3TL3UmSH3ShIIla2bLJho7SE9Jous6xMWUih
