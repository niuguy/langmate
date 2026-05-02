import AppKit
import Foundation

struct PreviewRequest: Decodable {
    let original: String
    let replacement: String
    let model: String?
}

struct PreviewResponse: Encodable {
    let action: String
    let text: String
}

final class PreviewController: NSObject, NSApplicationDelegate, NSWindowDelegate {
    private let request: PreviewRequest
    private var window: NSWindow!
    private var replacementView: NSTextView!
    private var didFinish = false

    init(request: PreviewRequest) {
        self.request = request
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApp.setActivationPolicy(.regular)

        let contentRect = NSRect(x: 0, y: 0, width: 600, height: 460)
        window = NSWindow(
            contentRect: contentRect,
            styleMask: [.titled, .closable, .resizable],
            backing: .buffered,
            defer: false
        )
        window.title = "LangMate Preview"
        window.center()
        window.level = .floating
        window.isReleasedWhenClosed = false
        window.collectionBehavior = [.canJoinAllSpaces, .fullScreenAuxiliary]
        window.delegate = self
        window.minSize = NSSize(width: 500, height: 360)

        let root = NSStackView()
        root.orientation = .vertical
        root.spacing = 8
        root.edgeInsets = NSEdgeInsets(top: 14, left: 14, bottom: 14, right: 14)
        root.translatesAutoresizingMaskIntoConstraints = false

        let header = NSStackView()
        header.orientation = .horizontal
        header.alignment = .centerY
        header.spacing = 8

        let titleLabel = label("Preview")
        let modelLabel = secondaryLabel(request.model ?? "Unknown model")
        let headerSpacer = NSView()
        headerSpacer.setContentHuggingPriority(.defaultLow, for: .horizontal)

        header.addArrangedSubview(titleLabel)
        header.addArrangedSubview(headerSpacer)
        header.addArrangedSubview(modelLabel)

        let originalLabel = secondaryLabel("Original")
        let originalScroll = textScrollView(text: request.original, editable: false)
        originalScroll.heightAnchor.constraint(equalToConstant: 86).isActive = true

        let replacementLabel = secondaryLabel("Replacement")
        let replacementScroll = textScrollView(text: request.replacement, editable: true)
        replacementView = replacementScroll.documentView as? NSTextView
        replacementScroll.heightAnchor.constraint(greaterThanOrEqualToConstant: 190).isActive = true

        let buttons = NSStackView()
        buttons.orientation = .horizontal
        buttons.alignment = .centerY
        buttons.spacing = 10

        let spacer = NSView()
        spacer.setContentHuggingPriority(.defaultLow, for: .horizontal)

        let cancelButton = NSButton(title: "Cancel", target: self, action: #selector(cancel))
        cancelButton.keyEquivalent = "\u{1b}"

        let copyButton = NSButton(title: "Copy", target: self, action: #selector(copyText))

        let replaceButton = NSButton(title: "Replace", target: self, action: #selector(replace))
        replaceButton.bezelStyle = .rounded
        replaceButton.keyEquivalent = "\r"
        replaceButton.keyEquivalentModifierMask = [.command]

        buttons.addArrangedSubview(spacer)
        buttons.addArrangedSubview(cancelButton)
        buttons.addArrangedSubview(copyButton)
        buttons.addArrangedSubview(replaceButton)

        root.addArrangedSubview(header)
        root.addArrangedSubview(originalLabel)
        root.addArrangedSubview(originalScroll)
        root.addArrangedSubview(replacementLabel)
        root.addArrangedSubview(replacementScroll)
        root.addArrangedSubview(buttons)

        window.contentView = NSView()
        window.contentView?.addSubview(root)

        NSLayoutConstraint.activate([
            root.leadingAnchor.constraint(equalTo: window.contentView!.leadingAnchor),
            root.trailingAnchor.constraint(equalTo: window.contentView!.trailingAnchor),
            root.topAnchor.constraint(equalTo: window.contentView!.topAnchor),
            root.bottomAnchor.constraint(equalTo: window.contentView!.bottomAnchor),
        ])

        window.makeKeyAndOrderFront(nil)
        NSRunningApplication.current.activate(options: [.activateAllWindows])
        window.makeFirstResponder(replacementView)
    }

    func windowWillClose(_ notification: Notification) {
        if !didFinish {
            finish(action: "cancel", text: request.replacement)
        }
    }

    @objc private func replace() {
        finish(action: "replace", text: replacementView.string)
    }

    @objc private func copyText() {
        finish(action: "copy", text: replacementView.string)
    }

    @objc private func cancel() {
        finish(action: "cancel", text: request.replacement)
    }

    private func finish(action: String, text: String) {
        if didFinish {
            return
        }
        didFinish = true
        window.delegate = nil

        let response = PreviewResponse(action: action, text: text)
        do {
            let data = try JSONEncoder().encode(response)
            FileHandle.standardOutput.write(data)
            FileHandle.standardOutput.write(Data("\n".utf8))
        } catch {
            fputs("Failed to encode response: \(error)\n", stderr)
        }
        NSApp.terminate(nil)
    }

    private func label(_ title: String) -> NSTextField {
        let field = NSTextField(labelWithString: title)
        field.font = NSFont.boldSystemFont(ofSize: 13)
        return field
    }

    private func secondaryLabel(_ title: String) -> NSTextField {
        let field = NSTextField(labelWithString: title)
        field.font = NSFont.systemFont(ofSize: 12)
        field.textColor = .secondaryLabelColor
        return field
    }

    private func textScrollView(text: String, editable: Bool) -> NSScrollView {
        let scrollView = NSScrollView()
        scrollView.borderType = .bezelBorder
        scrollView.hasVerticalScroller = true
        scrollView.translatesAutoresizingMaskIntoConstraints = false

        let textView = NSTextView()
        textView.string = text
        textView.isEditable = editable
        textView.isSelectable = true
        textView.font = NSFont.systemFont(ofSize: 13)
        textView.isRichText = false
        textView.allowsUndo = editable
        textView.textContainerInset = NSSize(width: 7, height: 7)
        textView.autoresizingMask = [.width]
        textView.textContainer?.widthTracksTextView = true

        scrollView.documentView = textView
        return scrollView
    }
}

let input = FileHandle.standardInput.readDataToEndOfFile()

do {
    let request = try JSONDecoder().decode(PreviewRequest.self, from: input)
    let app = NSApplication.shared
    let controller = PreviewController(request: request)
    app.delegate = controller
    app.run()
} catch {
    fputs("Failed to decode request: \(error)\n", stderr)
    exit(1)
}
