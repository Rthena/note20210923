package redis;

import java.util.ArrayList;
import java.util.List;

public class DistributeLock {


    public boolean tryLock_with_lua(String key, String UniqueId, int seconds) {
        String lua_scripts = "if redis.call('setnx',KEYS[1],ARGV[1]) == 1 then" +
                "redis.call('expire',KEYS[1],ARGV[2]) return 1 else return 0 end";
        List<String> keys = new ArrayList<>();
        List<String> values = new ArrayList<>();
        keys.add(key);
        values.add(UniqueId);
        values.add(String.valueOf(seconds));
//        Object result = jedis.eval(lua_scripts, keys, values);
        //判断是否成功
        return false;//result.equals(1L);
    }
}
