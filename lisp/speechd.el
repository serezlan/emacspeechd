(defsubst dtk-interp-set-voice (module-name voice-name)
  (cl-declare (special dtk-speaker-process))
  (process-send-string dtk-speaker-process (format "tts_set_voice %s %s\n" module-name voice-name)))

(defsubst dtk-interp-set-volume (volume)
  (cl-declare (special dtk-speaker-process))
  (process-send-string dtk-speaker-process (format "tts_set_volume %s\n" volume)))

(defun dtk-set-voice (module-name voice-name)
  "Set speech dispatcher voice."
  (cl-declare (special dtk-speaker-process
                       dtk-speak-server-initialized
                       dtk-quiet))
  (unless dtk-quiet
    (when dtk-speak-server-initialized
      (dtk-interp-set-voice module-name voice-name))))

(defun dtk-set-volume (volume)
  "Set speech dispatcher volume."
  (cl-declare (special dtk-speaker-process
                       dtk-speak-server-initialized
                       dtk-quiet))
  (unless dtk-quiet
    (when dtk-speak-server-initialized
      (dtk-interp-set-volume volume))))
