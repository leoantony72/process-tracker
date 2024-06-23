document.addEventListener("DOMContentLoaded", () => {
  const listElement = document.getElementById("dynamic-list");
  const criticalProcesses = [
    "smss.exe",
    "wininit.exe",
    "services.exe",
    "lsass.exe",
    "winlogon.exe",
    "csrss.exe",
    "spoolsv.exe",
    "Igfx EP Module.exe",
    "Igfx EC Module.exe",
    "vmms.exe",
    "atkexComSvc.exe",
    "SearchProtocolHost.exe",
    "MpDefenderCoreService.exe",
    "SecurityHealthService.exe",
    "vmcompute.exe",
    "StartMenuExperienceHost.exe",
    "igfxHK.exe",
    "dllhost.exe",
    "sihost.exe",
    "vmwp.exe",
    "ShellExperienceHost.exe",
    "SettingSyncHost.exe",
    "ShellExperienceHost.exe",
    "smartscreen.exe",
    "conhost.exe",
    "WUDFHost.exe",
    "sppsvc.exe",
    "RtkNGUI64.exe",
    "TiWorker.exe",
    "TrustedInstaller.exe",
    "SystemSettingsBroker.exe",
    "MsMpEng.exe",
    "CompPkgSrv.exe",
    "ApplicationFrameHost.exe",
    "mmc.exe",
    "WmiPrvSE.exe",
    "dasHost.exe",
    "SgrmBroker.exe",
    "NisSrv.exe",
    "SgrmBroker.exe",
    "igfxCUIService.exe",
    "tasklist.exe",
    "SearchFilterHost.exe",
    "igfxEM.exe",
    "gameinputsvc.exe",
    "TextInputHost.exe",
    "wslhost.exe",
    "LsaIso.exe",
    "fontdrvhost.exe",
    "SearchIndexer.exe",
    "Taskmgr.exe",
    "RuntimeBroker.exe",
    "explorer.exe",
    "svchost.exe",
    "dwm.exe",
    "taskhostw.exe",
    "firewall.exe",
    "msmpeng.exe",
    "audiodg.exe",
    "ctfmon.exe",
    "taskmgr.exe",
    "msiexec.exe",
    "searchui.exe",
  ];

  // Fetch the list from a URL
  fetch("http://192.168.18.11/send-command?cmd=get&pid_name=blablaf", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
  fetch("http://192.168.18.11:900/getdata")
    .then((response) => response.json())
    .then((data) => {
      console.log(data.data.data[0]);
      const items = data.data.data; // Adjust based on the JSON structure
      items.forEach((item) => {
        if (!criticalProcesses.includes(item)) {
          const li = document.createElement("li");
          li.textContent = item;
          li.addEventListener("click", () => handleClick(item));
          listElement.appendChild(li);
        }
      });
    })
    .catch((error) => console.error("Error fetching data:", error));
});

function handleClick(item) {
  const userConfirmed = confirm(`Do you want to send a request for ${item}?`);
  if (userConfirmed) {
    fetch(`http://192.168.18.11:900/send-command?cmd=kill&pid_name=${item}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => response.json())
      .then((data) => console.log("Response:", data))
      .catch((error) => console.error("Error sending data:", error));
  }
  // Send a POST request when an item is clicked
}
