# Powershell script to install Mcutie

# Run as admin
If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
		Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $PSCommandArgs" -WorkingDirectory $pwd -Verb RunAs
		Exit
	}

$ProgramName = "mcutie"
$ConfigFileName = "config.yaml"
$InstallPath = "c:\program files\Mcutie"
$ProgramNameEXE = "mcutie.exe"

Write-Host "================ Instal Mcutie ================`n`n"
Write-Host " - stop Mcutie if it's running"
Get-Process $ProgramName -ErrorAction SilentlyContinue | Stop-Process -PassThru

Write-Host " - delete existing install directory"
Remove-Item $InstallPath -Recurse -ErrorAction SilentlyContinue

Write-Host " - Create install directory: $InstallPath"
new-item $InstallPath -itemtype directory  | Out-Null

Write-Host " - Install $ProgramName to $InstallPath"
Copy-Item "$PSScriptRoot\bin\win64\mcutie.exe" -Destination $InstallPath

Write-Host " - Install $ConfigFileName to $InstallPath"
Copy-Item "$PSScriptRoot\$ConfigFileName" -Destination $InstallPath

Write-Host " - run at login for all users"
$startup=[Environment]::GetFolderPath("CommonStartup")

$RunScript= "Test.ps1"
$ShCutLnk = "Mcutie.lnk"

#create shortcut
$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$startup\$ShCutLnk")
$Shortcut.TargetPath = "$InstallPath\$ProgramNameEXE"
$Shortcut.WorkingDirectory = "$InstallPath"
$Shortcut.Save()

Write-Host "Press return to quit"
Read-Host