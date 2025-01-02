SELECT c.id,c.name,c.end_point ,c.api_key ,cm.model_key ,cm.total_count ,cm.error_count ,cm.weight
FROM channels c
         left join channel_models cm
                   on c.id = cm.channel_id ;

SELECT c.id,c.name,c.end_point ,c.api_key ,IFNULL(count(cm.id),0) as model_count
FROM channels c
         left join channel_models cm
                   on c.id = cm.channel_id
group by c.id;