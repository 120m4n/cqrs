DROP TABLE IF EXISTS feeds;

CREATE TABLE feeds (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- add sample data
INSERT INTO feeds (url, title, description) VALUES
    ('http://feeds.feedburner.com/techcrunch', 'TechCrunch', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/ommalik', 'GigaOM', 'GigaOM is a blog-related media company that offers news, analysis, and opinions on startup companies, emerging technologies, and other technology related topics.'),
    ('http://feeds.feedburner.com/venturebeat', 'VentureBeat', 'VentureBeat is a news site focused on information about innovation for forward-thinking executives.'),
    ('http://feeds.feedburner.com/thenextweb', 'The Next Web', 'The Next Web is one of the world’s largest online publications that delivers an international perspective on the latest news about Internet technology, business and culture.'),
    ('http://feeds.feedburner.com/techcrunch/startups', 'TechCrunch Startups', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/social', 'TechCrunch Social', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/gaming', 'TechCrunch Gaming', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/mobile', 'TechCrunch Mobile', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/enterprise', 'TechCrunch Enterprise', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/fundings-exits', 'TechCrunch Fundings & Exits', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.'),
    ('http://feeds.feedburner.com/techcrunch/europe', 'TechCrunch Europe', 'TechCrunch is a group-edited blog that profiles the companies, products and events defining and transforming the new web.');