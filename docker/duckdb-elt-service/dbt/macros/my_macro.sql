{% macro hello(name) %}
    select 'Hello {{ name }}' as greeting
{% endmacro %}

