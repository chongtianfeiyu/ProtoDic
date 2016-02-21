//此文件由ProtoDic批量生成，不要编辑此文件
package proto{import flash.utils.Dictionary;import handler.*;import proto.*;public class ProtoDic{
public static var protoMap:Dictionary;
public static var classMap:Dictionary;
public static var handleMap:Dictionary;
public function ProtoDic(){
		protoMap = new Dictionary();
		classMap = new Dictionary();
		handleMap = new Dictionary();

		protoMap[222] = StudentInfo;//学生
		classMap[StudentInfo] = 222;
		handleMap[222] = new StudentInfo_handler();
		protoMap[998] = LoginS;//登录
		classMap[LoginS] = 998;
		handleMap[998] = new LoginS_handler();
		
}}}
