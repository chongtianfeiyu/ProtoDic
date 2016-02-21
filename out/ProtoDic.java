//此文件由ProtoDic批量生成，不要编辑此文件
package proto;
import java.util.HashMap;
import handler.*;
import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
public final class ProtoDic {
	public static HashMap<Integer, Object> protoMap;
	public static HashMap<Object, Integer> classMap;
	public static HashMap<Integer, Object> handleMap;
	public static void init(){
		protoMap = new HashMap<Integer, Object>();
		classMap = new HashMap<Object, Integer>();
		handleMap = new HashMap<Integer, Object>();

		//协议号与协议的对应
		protoMap.put(222, Proto.StudentInfo.class);//学生
		classMap.put(Proto.StudentInfo.class, 222);
		protoMap.put(998, Proto.LoginS.class);//登录
		classMap.put(Proto.LoginS.class, 998);
		
	}
	public static Proto.StudentInfo toStudentInfo(ByteString data) {Proto.StudentInfo.Builder b = Proto.StudentInfo.newBuilder();try {b.mergeFrom(data);}catch (InvalidProtocolBufferException e) {e.printStackTrace();}return b.build();}
	public static Proto.LoginS toLoginS(ByteString data) {Proto.LoginS.Builder b = Proto.LoginS.newBuilder();try {b.mergeFrom(data);}catch (InvalidProtocolBufferException e) {e.printStackTrace();}return b.build();}
	
}
