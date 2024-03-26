CREATE TABLE search_result (
    id SERIAL PRIMARY KEY,
    keyword TEXT,
    html_content TEXT,
    adword_amount INT,
    total_search_result VARCHAR(200),
    user_id VARCHAR(200),
    link_amount INT
);