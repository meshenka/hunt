--
-- Company
--

-- name: Companies :many
SELECT * FROM company
WHERE author_id=$1
ORDER BY created_at DESC;

-- name: CompanyByID :one
SELECT * FROM company
WHERE author_id=$1
AND id=$2;

--
-- Opportunity
--

-- name: Opportunities :many
SELECT o.* FROM opportunity AS o
INNER JOIN company AS c ON o.company_id = c.id
INNER JOIN opportunity_note AS n ON o.opportunity_id = n.id
WHERE o.author_id=$1
ORDER BY o.created_at DESC;

-- name: OpportunityByID :one
SELECT * FROM opportunity
WHERE author_id=$1
AND id=$2;

--
-- Notes
--

-- name: OpportunityNotes :many
SELECT * FROM opportunity_note
WHERE author_id=$1
AND opportunity_id=$2
ORDER BY created_at DESC;
