--
-- PostgreSQL database dump
--

-- Dumped from database version 14.9 (Ubuntu 14.9-1.pgdg20.04+1)
-- Dumped by pg_dump version 14.9 (Ubuntu 14.9-1.pgdg20.04+1)

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

--
-- Name: stat; Type: TYPE; Schema: public; Owner: ameps
--

CREATE TYPE public.stat AS ENUM (
    'en attente',
    'affecté',
    'résolu',
    'archivé'
);


ALTER TYPE public.stat OWNER TO ameps;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: info; Type: TABLE; Schema: public; Owner: ameps
--

CREATE TABLE public.info (
    id integer NOT NULL,
    agent character varying NOT NULL,
    material character varying NOT NULL,
    priority integer NOT NULL,
    rte character varying,
    detail text NOT NULL,
    estimate character varying,
    brips character varying,
    oups character varying,
    ameps character varying,
    ais character varying,
    source_id integer NOT NULL,
    created date NOT NULL,
    updated date,
    status character varying,
    event character varying NOT NULL,
    target character varying,
    doneby character varying,
    pilote character varying,
    action_date character varying,
    day_done character varying
);


ALTER TABLE public.info OWNER TO ameps;

--
-- Name: info_info_id_seq; Type: SEQUENCE; Schema: public; Owner: ameps
--

CREATE SEQUENCE public.info_info_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.info_info_id_seq OWNER TO ameps;

--
-- Name: info_info_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ameps
--

ALTER SEQUENCE public.info_info_id_seq OWNED BY public.info.id;


--
-- Name: source; Type: TABLE; Schema: public; Owner: ameps
--

CREATE TABLE public.source (
    id integer NOT NULL,
    name character varying NOT NULL,
    created date NOT NULL,
    code_gmao character varying
);


ALTER TABLE public.source OWNER TO ameps;

--
-- Name: source_id_seq; Type: SEQUENCE; Schema: public; Owner: ameps
--

CREATE SEQUENCE public.source_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.source_id_seq OWNER TO ameps;

--
-- Name: source_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ameps
--

ALTER SEQUENCE public.source_id_seq OWNED BY public.source.id;


--
-- Name: info id; Type: DEFAULT; Schema: public; Owner: ameps
--

ALTER TABLE ONLY public.info ALTER COLUMN id SET DEFAULT nextval('public.info_info_id_seq'::regclass);


--
-- Name: source id; Type: DEFAULT; Schema: public; Owner: ameps
--

ALTER TABLE ONLY public.source ALTER COLUMN id SET DEFAULT nextval('public.source_id_seq'::regclass);


--
-- Data for Name: info; Type: TABLE DATA; Schema: public; Owner: ameps
--

COPY public.info (id, agent, material, priority, rte, detail, estimate, brips, oups, ameps, ais, source_id, created, updated, status, event, target, doneby, pilote, action_date, day_done) FROM stdin;
237	KALER	UA	1		DJ REDRESSEUR 48V DECLENCHE				OUI		57	2023-04-26	\N	en attente	MP			\N	\N	\N
258	DANTAS	AT- 113- RLT	2		RLT à nettoyer PH 2 et PH 10 				OUI		60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
238	DANTAS	BATIMENT- infiltration et fissure mur	2		Salle de reunion et rame AB infiltration et fissure dans le mur. 		oui				60	2023-05-15	2023-05-15	en attente	Visite technique 			\N	\N	\N
239	DANTAS	Galerie HTA - infiltration + éclairage 	2		Infiltration d'eau + éclairage à reprendre dans sa globalité.		oui				60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
240	DANTAS	Galerie HTA - câble HTA non sécurisé	1		câbles HTA non mis en court-circuit 						60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
241	DANTAS	Galerie HTA - Escalier à néttoyer 	3		Escalier à nettoyer 				OUI		60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
242	DANTAS	BATIMENT- chauffages	2		Chauffages à remplacer dans l'ensemble du poste.		oui				60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
243	DANTAS	Extérieur- Barbelé	2		Barbelé à remettre.				OUI		60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
244	DANTAS	Extérieur- Eclairage  	2		Eclairage extérieur à reprendre en totalité.				OUI		60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
245	DANTAS	Extérieur- Désherbage 	1		désherbage à prévoir dans le poste.				OUI		60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
246	DANTAS	Extérieur- dalles	1		dalles cassées dans le poste.				OUI		60	2023-05-15	\N	en attente	Visite technique			\N	\N	\N
259	DANTAS	AT- 113 - Tresses de terre	2		Tresses de terre à isoler. 				OUI		60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
247	DANTAS	local ventil. tr612	2		Probléme de sérrure 				OUI		60	2023-05-15	2023-05-16	en attente	Visite technique			\N	\N	\N
248	DANTAS	Loge TR 611/612 - sécurisation échelle	1		chainette manquante au niveau de échelle extérieur.				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
249	DANTAS	Stockage 	2		tourets de câbles stockés à coté de la loge TR gênant l'accés. 				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
250	DANTAS	WC	2		CÂBLE ÉLECTRIQUE ISOLÉ PAR DU SCOTCH \r\n\r\nA dépossé ? 				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
251	DANTAS	Rames- BPAG	1		Rames GH-AB-CD-EF : il n'y a pas de BPAG.				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
252	DANTAS	Local Batteries	2		il n'y a pas de point d'eau ou rince oeil.\r\n				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
253	DANTAS	Matériaux à récupérer 	2		TC et EPAMI et redresseur dans les salles rames et couloir à récupérer 				OUI		60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
255	DANTAS	Contrôle commande- PC C3S 	1		écran ne fonctionne pas. 				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
257	DANTAS	TR 612- Fuite	2		légère fuite sur le transfo à contrôler. 				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
273	DA SILVA	BATIMENT - Caillebotis et Faux plancher	3		Caillebotis gondolés dans chemin gradins, AT et salle contrôle commande				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
260	DANTAS	Rame- I - Vibration	2		Forte vibration de la tôle. 				OUI		60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
262	DA SILVA	GALERIE HTA - Terre	2		Terre manquant sur une grande longueur des chemins de câbles.				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
266	DA SILVA	GALERIE HTA - Eclairage	3		Éclairage galerie manquant et HS côté TR 611 suite infiltration d'eau				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
267	DA SILVA	BATIMENT - Portes	3		Groom des portes HS:RAME A, C (Côté couloir + sortie extérieur + SS),G (SS),Salle réserve, TR 613				OUI		55	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
269	DA SILVA	SALLES RAMES - Bloc de secours	2		Toutes les salles rames n'ont pas de bloc de secours		oui		OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
270	DA SILVA	BATIMENT - Bloc de secours	2		Blocs de secours sans autocollant de direction dans le couloirs				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
271	DA SILVA	PALAN - Condos	3		Palan non révisé depuis 2015				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
272	DA SILVA	RAME B, G, H - Viserie	3		Manque boulon pour fermer les capots de racks des cellules 				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
274	DA SILVA	LOGE AT 112 - Ventilation	2		Ventilation extraction HS de la loge AT 112				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
268	DA SILVA	CABLES HTA - Sous-sol	3		Non repérage des câbles HTA des départs.A voir avec l'exploitation						55	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
265	DA SILVA	LOGE TRANSFO - Portes	3		Portes des loges TR 612 + TR 613 HS. Très difficile voir impossible à ouvrir 						55	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
275	DA SILVA	PSEM - Armoire TR 612	3		Porte armoire TR 612 HS. Charnières cassées				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
278	DA SILVA	POSTE SOURCE - Dalles	2		Dalles cassées piste côté grille TR 612				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
277	DA SILVA	POSTE SOURCE - Fuite à la terre	1		Fuite à la terre courant alternatif côté TR 613. Environ 190V dans le neutre				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
276	DA SILVA	BATIMENT - Stockage	3		Stockage matériel non rangé correctement et non balisé dans Salle atelier et dans la loge TR 612			 	OUI		55	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
279	DA SILVA	SALLE CC - Faux plancher	1		Faux plancher a reprendre. Risque de chute				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
280	DA SILVA	SALLE REPOS - Hygiène	3		Salle de repos non propre et occupé par d'autres sociétés  				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
281	DA SILVA	POSTE SOURCE - Désherbage	3		Désherbage est à prévoir. Empêchement d'accéder à la galerie et ouvrir la porte couloir AT				OUI		55	2023-05-16	2024-02-15	résolu	Visite technique	15/02/2024		\N	\N	\N
264	DA SILVA	TR 612 - Silicagel	3		Silicagel à changer				OUI		55	2023-05-16	2024-02-15	résolu	Visite technique	10/02/2024		\N	\N	\N
256	DANTAS	Aéros -Tr 611/612	2		Prévoir nettoyage sur les aéros. 				OUI		60	2023-05-16	2024-02-15	résolu	Visite technique	08/02/2024		\N	\N	\N
261	DA SILVA	BPAG	1		Toute l'extension bâtiment, galeries HTA, TR612 + TR613 possède des BPAG HS. Soit fil en l'air soit ne remonte pas soit non alarmé 				OUI		55	2023-05-16	2024-02-15	résolu	Visite technique	22/01/2024		\N	\N	\N
282	DA SILVA	CUVE BARBOTAGE - Fuite	2		Fuite dans les cuves de barbotages des AT. Une cuve est totalement sèche. En cas de feu dans l'AT l'huile en ébullition ne sera pas refroidi				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
283	DA SILVA	TRAPPES ACCES - Fuite	2		Toutes les trappes d'accès aux galeries ne sont plus étanches				OUI		55	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
286	SALABERT	POSTE SOURCE - Eclairage	2		Plusieurs éclairage ne fonctionne pas.				OUI		45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
287	SALABERT	BATIMENT - Stockage	2		Matériel stocké à plusieurs endroit du poste sans balisage.				oui		45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
288	SALABERT	TR 613 - Cuve de barbotage	2		Manque vérin pour l'ouverture de la trappe de la cuve de barbotage.				OUI		45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
309	DANTAS	LOGE RPN + CONSERVATEUR - BPAG	1		pas de BPAG dans les loges des conservateur et des RPN.						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
292	DANTAS	Rame GH - Disj HTA	2		Disj HTA  à contrôler à 1250 A				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
293	DANTAS	Rame CD - Disj HTA	2		Disj HTA à contôler à 630A 				OUI		60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
294	DANTAS	Disj HTA - Entretien	2		contrôle et entretien disj HTA FPR ( 3 disj de 630 A, 1 de 1250A et de 400A).				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
295	DANTAS	Cuisine 	3		Enlèvement gazinière et frigo.				OUI		60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
296	DANTAS	Salle de commande - climatisation 	2		climatisation ne fonctionne pas.						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
297	DANTAS	Salle de réunion - vitre	3		vitre fissurée. 						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
299	DANTAS	Extincteurs 	2		Plusieurs Extincteur non contrôlés depuis le 12/2014.						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
300	DANTAS	Salle de réunion - mobilier 	3		mobilier à changer trop ancien. 						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
301	DANTAS	GALERIE HTA - boite de dérivation 	2		boite de dérivation cassée. 						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
302	DANTAS	téléphone du poste 	2		téléphone du poste qui ne fonctionne pas. 						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
303	DANTAS	Départ R91 -Trémie 	2		Trémie à reboucher  départ R91 						60	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
310	DANTAS	RPN -TR 612 inondée	2		inondations 						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
311	DANTAS	RPN -TR 611 déconnectée 	2		RPN TR 611 déconnectée. 						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
304	DANTAS	LOGE TRANSFO- Éclairage + interrupteur 	1		Éclairages mal positionné et hors service dans les loges des TR +interrupteur mal fixés. 						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
305	DANTAS	CI 	2		les connexions dans les CI sont oxydé en plus d'un début de  corrosion sur les contactes. 						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
307	DANTAS	GRILLES TR - sabots de terre	2		sabots de terre des grilles mal positionnés.						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
313	DANTAS	LOGE CONSERVATEUR - porte 	2		poignée de porte cassé et démonté d'une des loge des conservateur au N+1.						45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
312	DANTAS	AÉROS- filtres	2		préfiltres des aéros  impossible à démonter au N+1.\r\n						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
290	SALABERT	SOUS-SOL anti-panique	1		Dans le sous-sol du nouveau poste une barre anti-panique a été cassée.				OUI		45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
291	SALABERT	SOUS-SOL balisage	1		Balisage insufisant devant des caniveaux ouvert.				OUI		45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
289	SALABERT	SOUS-SOL déjection	3		Déjection animal retrouvée dans le sous-sol.				OUI		45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
254	DANTAS	Contrôle commande- EMIC 	1		 EMIC ne fonctionne pas. 				OUI		60	2023-05-16	2023-07-31	résolu	Visite technique	01/06/2023	REDON	\N	\N	\N
314	DANTAS	SSI	2		Boîtier SSI cassé au N+1. \r\n						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
308	DANTAS	TRÉMIES 	1		Infiltration d'eau car les trémies ne sont pas bouchée N-0. Trémies non bouchée entre loge TR et blocs aéroréfrigérants N + 1. trémies à reboucher dans la salle des TSA au N+1.						45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
315	DANTAS	MISE A LA TERRE 	1		Dans le local TCFM le faux planché métallique n'est pas à la terre. 						45	2023-05-16	\N	en attente	Visite technique			\N	\N	\N
285	SALABERT	 BPAG ET PRISE DE COURANT	2		Fixation des BPAG et des Prise de courant a revoir dans le nouveau bâtiment.				OUI		45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
298	DANTAS	POSTE SOURCE - Eclairage	1		Éclairage intérieur et extérieur à reprendre dans sa globalité. 						60	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
306	DANTAS	EXCREMENTS D'OISEAUX	1		excréments d’oiseaux un peut partout au poste. 		oui				45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
284	SALABERT	POSTE SOURCE - passerelle	1		Les passerelles métallique des coursives et les rampes d'escalier ne sont pas à la terre.\r\n		OUI				45	2023-05-16	2023-05-16	en attente	Visite technique			\N	\N	\N
263	DA SILVA	TR 612 - Aéros	3		Bruit anormal sur un moteur aéro bas droite				OUI		55	2023-05-16	2023-10-09	résolu	Visite technique			\N	\N	\N
317	Rodriguez	Pa	1		PA sur redresseur secours. Pb alim				Oui		48	2024-04-04	\N	affecté	Panne			\N	\N	\N
\.


--
-- Data for Name: source; Type: TABLE DATA; Schema: public; Owner: ameps
--

COPY public.source (id, name, created, code_gmao) FROM stdin;
53	Alsace	2023-01-10	ALSAC
55	Ampère	2023-01-10	AMPER
60	Argenteuil	2023-01-17	ARG6
58	Billancourt	2023-01-14	BILLA
57	Boule	2023-01-10	BOULE
34	Courbevoie	2022-11-21	CZBEV
28	Danton	2022-10-31	DANT5
41	La Briche	2022-11-23	BRICH
42	Levallois	2022-11-23	LEVAL
43	Menus	2022-11-23	MENUS
44	Nanterre	2022-11-23	NANTE
20	Novion	2022-10-17	NOVIO
45	Puteaux	2022-11-23	PUTEA
46	Rueil	2022-11-23	RUEIL
47	St Ouen	2022-11-23	SSOU5
48	Tilliers	2022-11-23	TILLI
\.


--
-- Name: info_info_id_seq; Type: SEQUENCE SET; Schema: public; Owner: ameps
--

SELECT pg_catalog.setval('public.info_info_id_seq', 317, true);


--
-- Name: source_id_seq; Type: SEQUENCE SET; Schema: public; Owner: ameps
--

SELECT pg_catalog.setval('public.source_id_seq', 64, true);


--
-- Name: info infos_pkey; Type: CONSTRAINT; Schema: public; Owner: ameps
--

ALTER TABLE ONLY public.info
    ADD CONSTRAINT infos_pkey PRIMARY KEY (id);


--
-- Name: source sources_pkey; Type: CONSTRAINT; Schema: public; Owner: ameps
--

ALTER TABLE ONLY public.source
    ADD CONSTRAINT sources_pkey PRIMARY KEY (id);


--
-- Name: info fk_source; Type: FK CONSTRAINT; Schema: public; Owner: ameps
--

ALTER TABLE ONLY public.info
    ADD CONSTRAINT fk_source FOREIGN KEY (source_id) REFERENCES public.source(id);


--
-- Name: TABLE info; Type: ACL; Schema: public; Owner: ameps
--

GRANT ALL ON TABLE public.info TO web;


--
-- Name: SEQUENCE info_info_id_seq; Type: ACL; Schema: public; Owner: ameps
--

GRANT SELECT,USAGE ON SEQUENCE public.info_info_id_seq TO web;


--
-- Name: TABLE source; Type: ACL; Schema: public; Owner: ameps
--

GRANT ALL ON TABLE public.source TO web;


--
-- Name: SEQUENCE source_id_seq; Type: ACL; Schema: public; Owner: ameps
--

GRANT SELECT,USAGE ON SEQUENCE public.source_id_seq TO web;


--
-- PostgreSQL database dump complete
--

