package logic

const (
	//#1.0.1 包含時間加權
	selectOrderJoinWorkorderGroupByManorderIncludeCalTime = `select orders.order_datetime,orders.due_datetime,T.manorder_id,(orders.order_info->>'order_id')::text as order_id,T.qty as total_qty,orders.qty,T.required_time,T.worked_time,COALESCE(T.acc_good,0)as acc_good,COALESCE(T.acc_ng,0)as acc_ng from mes.orders left join
	(select manorder_id,COALESCE(sum(worked_time),0) as worked_time,COALESCE(sum(required_time),0) as required_time,COALESCE(sum(acc_good),0) as acc_good, COALESCE(sum(acc_ng),0) as acc_ng, sum(qty) as qty from
	(select (work_orders.order_info->>'manorder_id')::text as manorder_id, (work_orders.order_info->>'step_id')::int as step, (work_orders.product->'route'->(((work_orders.order_info->>'step_id')::int)-1)->'lines'->0->>'process_time')::int as process_time, ((COALESCE((work_orders.state->>'acc_good')::int,0))*(work_orders.product->'route'->(((work_orders.order_info->>'step_id')::int)-1)->'lines'->0->>'process_time')::int) as worked_time,((COALESCE((work_orders.order_info->>'qty')::int,0))*(work_orders.product->'route'->(((work_orders.order_info->>'step_id')::int)-1)->'lines'->0->>'process_time')::int) as required_time, (work_orders.assignment->>'duration')::int as duration, (work_orders.order_info->>'qty')::int as qty,COALESCE((work_orders.state->>'acc_good')::int,0) as acc_good,COALESCE((work_orders.state->>'acc_ng')::int,0) as acc_ng from mes.work_orders order by pk) as rTable	
	group by manorder_id) T
	on (orders.order_info->>'manorder_id')::text = T.manorder_id 
	order by T.manorder_id`
	//#1.0.1 只統計數量不加權
	selectOrderJoinWorkorderGroupByManorder = `select orders.order_datetime,orders.due_datetime,T.manorder_id,(orders.order_info->>'order_id')::text as order_id,T.qty as total_qty,orders.qty,COALESCE(T.acc_good,0)as acc_good,COALESCE(T.acc_ng,0)as acc_ng from mes.orders left join 
	(select manorder_id,COALESCE(sum(acc_good),0) as acc_good, COALESCE(sum(acc_ng),0) as acc_ng, sum(qty) as qty from
	(select (work_orders.order_info->>'manorder_id')::text as manorder_id, (work_orders.order_info->>'qty')::int as qty,(work_orders.state->>'acc_good')::int as acc_good,(work_orders.state->>'acc_ng')::int as acc_ng from mes.work_orders) as rTable
	group by manorder_id) T
	on (orders.order_info->>'manorder_id')::text = T.manorder_id
	order by T.manorder_id`
)
