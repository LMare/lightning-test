select
    id,
    upper(message) as message_upper
from {{ ref('stg_example') }}

