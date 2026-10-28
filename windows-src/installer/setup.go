//go:build windows

package main

import (
    "embed"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "syscall"
    "unsafe"
)

//go:embed payload/HOI4-Focus-Studio.exe payload/icon.ico payload/README.txt
var payload embed.FS

const (
    productName = "HOI4 Focus Studio"
    regKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\HOI4FocusStudio`
    mbYesNo = 0x04
    mbIconQuestion = 0x20
    mbIconError = 0x10
    mbIconInfo = 0x40
    mbDefaultSecond = 0x100
    idYes = 6
)

func message(text string, flags uintptr) uintptr {
    title,_ := syscall.UTF16PtrFromString(productName)
    body,_ := syscall.UTF16PtrFromString(text)
    proc := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
    ret,_,_ := proc.Call(0,uintptr(unsafe.Pointer(body)),uintptr(unsafe.Pointer(title)),flags)
    return ret
}

func fatal(err error){message("Ошибка: "+err.Error(),mbIconError);os.Exit(1)}
func yesNo(question string) bool{return message(question,mbYesNo|mbIconQuestion|mbDefaultSecond)==idYes}

func installDir() string {
    local := os.Getenv("LOCALAPPDATA")
    if local == "" {home,_:=os.UserHomeDir();local=filepath.Join(home,"AppData","Local")}
    return filepath.Join(local,"Programs","HOI4 Focus Studio")
}
func savesDir() string {
    local := os.Getenv("LOCALAPPDATA")
    if local == "" {home,_:=os.UserHomeDir();local=filepath.Join(home,"AppData","Local")}
    return filepath.Join(local,"HOI4FocusStudio")
}
func startMenu() string {
    appdata := os.Getenv("APPDATA")
    if appdata == "" {home,_:=os.UserHomeDir();appdata=filepath.Join(home,"AppData","Roaming")}
    return filepath.Join(appdata,"Microsoft","Windows","Start Menu","Programs","HOI4 Focus Studio.lnk")
}
func userDesktop() string {
    home,_:=os.UserHomeDir()
    return filepath.Join(home,"Desktop","HOI4 Focus Studio.lnk")
}
func writeEmbedded(name, location string) error {
    b,e := payload.ReadFile("payload/"+name)
    if e!=nil{return e}
    tmp:=location+".new"
    if err:=os.WriteFile(tmp,b,0644);err!=nil{return err}
    if _,statErr:=os.Stat(location);statErr==nil {if err:=os.Remove(location);err!=nil {_=os.Remove(tmp);return err}}
    if err:=os.Rename(tmp,location);err!=nil {_ = os.Remove(tmp);return err}
    return nil
}
func copySelf(path string) error {
    src,e:=os.Executable();if e!=nil{return e}
    in,e:=os.Open(src);if e!=nil{return e};defer in.Close()
    out,e:=os.Create(path+".new");if e!=nil{return e}
    _,e=io.Copy(out,in);if ce:=out.Close();e==nil{e=ce}
    if e!=nil{_ = os.Remove(path+".new");return e}
    if _,statErr:=os.Stat(path);statErr==nil {if err:=os.Remove(path);err!=nil {_=os.Remove(path+".new");return err}}
    if e=os.Rename(path+".new",path);e!=nil {_=os.Remove(path+".new")}
    return e
}
func psQuote(s string)string{return "'"+strings.ReplaceAll(s,"'","''")+"'"}
func shortcut(path, target, icon, workdir string) error {
    if err:=os.MkdirAll(filepath.Dir(path),0755);err!=nil{return err}
    script := "$w=New-Object -ComObject WScript.Shell;"+
        "$s=$w.CreateShortcut("+psQuote(path)+");"+
        "$s.TargetPath="+psQuote(target)+";"+
        "$s.IconLocation="+psQuote(icon)+";"+
        "$s.WorkingDirectory="+psQuote(workdir)+";"+
        "$s.Description='HOI4 Focus Studio — редактор национальных фокусов';"+
        "$s.Save()"
    cmd:=exec.Command("powershell.exe","-NoProfile","-NonInteractive","-ExecutionPolicy","Bypass","-Command",script)
    cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}
    out,e:=cmd.CombinedOutput()
    if e!=nil{return fmt.Errorf("создание ярлыка: %v (%s)",e,strings.TrimSpace(string(out)))}
    return nil
}
func registryAdd(key, val, valType, data string)error{
    args:=[]string{"add",key,"/v",val,"/t",valType,"/d",data,"/f"}
    cmd:=exec.Command("reg.exe",args...);cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}
    out,e:=cmd.CombinedOutput()
    if e!=nil{return fmt.Errorf("запись в реестр: %w (%s)",e,strings.TrimSpace(string(out)))}
    return nil
}
func register(dir string) error {
    pairs:=[][3]string{
       {"DisplayName","REG_SZ",productName},
       {"DisplayVersion","REG_SZ","1.1.0"},
       {"Publisher","REG_SZ","HOI4 Focus Studio"},
       {"InstallLocation","REG_SZ",dir},
       {"DisplayIcon","REG_SZ",filepath.Join(dir,"icon.ico")},
       {"UninstallString","REG_SZ",`"`+filepath.Join(dir,"Uninstall.exe")+`" /uninstall`},
       {"NoModify","REG_DWORD","1"},
       {"NoRepair","REG_DWORD","1"},
    }
    for _,p:=range pairs{if e:=registryAdd(regKey,p[0],p[1],p[2]);e!=nil{return e}}
    return nil
}
func install(){
    dir:=installDir()
    if !yesNo("Установить HOI4 Focus Studio для текущего пользователя?\n\nПапка: "+dir+"\n\nПриложение работает без интернета, но требует Microsoft Edge для собственного окна.\n\nАвтосохранения будут размещаться отдельно от программы."){return}
    if e:=os.MkdirAll(dir,0755);e!=nil{fatal(e)}
    for _,name:=range []string{"HOI4-Focus-Studio.exe","icon.ico","README.txt"} {
       if e:=writeEmbedded(name,filepath.Join(dir,name));e!=nil{fatal(fmt.Errorf("не удалось записать %s: %w\nЗакройте уже запущенное приложение и попробуйте снова",name,e))}
    }
    if e:=copySelf(filepath.Join(dir,"Uninstall.exe"));e!=nil{fatal(e)}
    app:=filepath.Join(dir,"HOI4-Focus-Studio.exe")
    ico:=filepath.Join(dir,"icon.ico")
    if e:=shortcut(startMenu(),app,ico,dir);e!=nil{fatal(e)}
    if yesNo("Добавить ярлык HOI4 Focus Studio на рабочий стол?"){
        if e:=shortcut(userDesktop(),app,ico,dir);e!=nil {
            message("Установка выполнена, но ярлык на рабочем столе создать не удалось:\n"+e.Error(),mbIconInfo)
        }
    }
    if e:=register(dir);e!=nil{fatal(e)}
    if yesNo("HOI4 Focus Studio установлен!\n\nПрограмма доступна через меню «Пуск».\n\nЗапустить редактор сейчас?"){
        cmd:=exec.Command(app)
        if e:=cmd.Start();e!=nil{fatal(e)}
        go func(){_ = cmd.Wait()}()
    }
}
func uninstall(){
    dir:=installDir()
    if !yesNo("Удалить HOI4 Focus Studio?\n\nБудут удалены программа и ярлыки.\nСохранения можно оставить на компьютере.\n\nПеред удалением закройте все окна приложения."){return}
    deleteSaves:=yesNo("Удалить ТАКЖЕ все проекты, версии и локальные сохранения?\n\nЕсли нажать «Нет», сохранения останутся и будут доступны после повторной установки.\nДля переноса на другой компьютер рекомендуем предварительно скачать JSON-копию.")
    _=os.Remove(startMenu());_ = os.Remove(userDesktop())
    cmd:=exec.Command("reg.exe","delete",regKey,"/f");cmd.SysProcAttr=&syscall.SysProcAttr{HideWindow:true};_ = cmd.Run()
    errs:=[]string{}
    for _,name:=range []string{"HOI4-Focus-Studio.exe","icon.ico","README.txt"} {
        if e:=os.Remove(filepath.Join(dir,name));e!=nil&&!os.IsNotExist(e){errs=append(errs,e.Error())}
    }
    if deleteSaves {if e:=os.RemoveAll(savesDir());e!=nil{errs=append(errs,"сохранения: "+e.Error())}}
    if len(errs)>0 {
        message("Часть файлов не удалось удалить. Убедитесь, что редактор закрыт.\n\n"+strings.Join(errs,"\n"),mbIconError)
        return
    }
    message("Удаление завершено.\n\n"+map[bool]string{true:"Сохранения удалены.",false:"Ваши сохранения сохранены в "+savesDir()}[deleteSaves],mbIconInfo)
    // The uninstaller cannot delete its own running executable. A hidden
    // PowerShell child removes it after this process exits.
    cleanup := "Start-Sleep -Seconds 3; " +
        "Remove-Item -LiteralPath "+psQuote(filepath.Join(dir,"Uninstall.exe"))+" -Force -ErrorAction SilentlyContinue; " +
        "Remove-Item -LiteralPath "+psQuote(dir)+" -Force -ErrorAction SilentlyContinue"
    child:=exec.Command("powershell.exe","-NoProfile","-NonInteractive","-WindowStyle","Hidden","-Command",cleanup)
    child.SysProcAttr=&syscall.SysProcAttr{HideWindow:true}
    _=child.Start()
}
func main(){
    exe,_:=os.Executable()
    if len(os.Args)>1&&os.Args[1]=="/uninstall"||strings.EqualFold(filepath.Base(exe),"Uninstall.exe"){uninstall();return}
    install()
}
