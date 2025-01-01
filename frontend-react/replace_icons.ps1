param(
    [string]$Path
)

$iconMappings = @{
    "IconHorse" = "MdPets";
    "IconPlus" = "FiPlus";
    "IconCircleCheck" = "MdCheckCircle";
    "IconMars" = "MdMale";
    "IconVenus" = "MdFemale";
    "IconSearch" = "FiSearch";
    "IconEdit" = "FiEdit";
    "IconTrash" = "FiTrash2"
}

$content = Get-Content $Path -Raw

# Remove Tabler import
$content = $content -replace "import \{[^}]*\} from '@tabler/icons-react';", ""

# Add React Icons import
$reactIconImports = @()
foreach ($oldIcon in $iconMappings.Keys) {
    if ($content -match $oldIcon) {
        $newIcon = $iconMappings[$oldIcon]
        $reactIconImports += "import { $newIcon } from 'react-icons/md';"
        $reactIconImports += "import { FiPlus, FiSearch, FiEdit, FiTrash2 } from 'react-icons/fi';"
        $content = $content -replace $oldIcon, $newIcon
    }
}

# Insert imports at the top
$importBlock = $reactIconImports -join "`n"
$content = $importBlock + "`n`n" + $content

Set-Content -Path $Path -Value $content
