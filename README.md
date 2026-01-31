# epub-renamer

## Scan files

Scans files.

### Command

```
epub-renamer scan (input-dir) (output-file)
```

* `input-dir` Directory containing input books.
* `output-file` Output CSV file which will be populated with metadata and new filenames. 

### Config file

`config.yml`
```
useOriginalDirectoryLayout: true
filenameTemplateParts:
  - '{{ .Metadata.Creator }} - {{ .Metadata.Title }}'
```

* `useOriginalDirectoryLayout` `bool` Whether to use original directory layout.
* `filenameTemplateParts` `string` Filename template parts.

#### Template variables

| Name                  | Description           |
|-----------------------|-----------------------|
| .Metadata.Title       | Metadata: Title       |
| .Metadata.Language    | Metadata: Language    |
| .Metadata.Identifier  | Metadata: Identifier  |
| .Metadata.Creator     | Metadata: Creator     |
| .Metadata.Contributor | Metadata: Contributor |
| .Metadata.Publisher   | Metadata: Publisher   |
| .Metadata.Subject     | Metadata: Subject     |
| .Metadata.Description | Metadata: Description |
| .Metadata.Event.Name  | Metadata: Event Name  |
| .Metadata.Event.Date  | Metadata: Event Date  |
| .Metadata.Type        | Metadata: Type        |
| .Metadata.Format      | Metadata: Format      |
| .Metadata.Source      | Metadata: Source      |
| .Metadata.Relation    | Metadata: Relation    |
| .Metadata.Coverage    | Metadata: Coverage    |
| .Metadata.Rights      | Metadata: Rights      |

## Rename files

Renames (copies) files according to the CSV file provided.

This sub-command uses the `input` and `output` columns of the provided CSV file.

```
epub-renamer rename (input-dir) (input-file) (output-dir)
```
